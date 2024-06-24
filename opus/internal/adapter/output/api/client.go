package api

import (
	"context"
	"fmt"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/api"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log/slog"
)

type Client struct {
	ytSvc *youtube.Service
}

func NewYouTubeClient(ctx context.Context, apiKey string) (*Client, error) {
	svc, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &Client{ytSvc: svc}, nil
}

func (c *Client) FetchVideoDetailsByVideoIDs(ctx context.Context, videoIDs []string) ([]*api.VideoDetail, error) {
	call := c.ytSvc.Videos.List([]string{"snippet", "contentDetails", "liveStreamingDetails"}).Id(videoIDs...)
	call = call.Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	vds := make([]*api.VideoDetail, 0, len(resp.Items))

	for _, i := range resp.Items {
		vd := api.NewVideoDetail(i.Id)

		if i.Snippet == nil {
			return nil, fmt.Errorf("snippet is not found")
		}

		// publishedAt
		pa, err := synchro.ParseISO[tz.AsiaTokyo](i.Snippet.PublishedAt)
		if err != nil {
			slog.Error(
				"failed to parse 'publishedAt' for video ID; "+i.Id,
				"publishedAt", i.Snippet.PublishedAt,
				slog.Group("fetchVideoDetailsByVideoIDs", "error", err),
			)
			return nil, err
		}
		vd.SetPublishedAtUnix(pa.Unix())

		// liveBroadcastContent
		switch i.Snippet.LiveBroadcastContent {
		case "live":
			vd.SetStatus(api.Live)
			vd.SetChatID(extractChatID(i.LiveStreamingDetails))
			sa, err := extractScheduledAtUnix(i.LiveStreamingDetails)
			if err != nil {
				return nil, err
			}
			if err := vd.SetScheduledAtUnix(sa); err != nil {
				return nil, err
			}
		case "upcoming":
			vd.SetStatus(api.Upcoming)
			vd.SetChatID(extractChatID(i.LiveStreamingDetails))
			sa, err := extractScheduledAtUnix(i.LiveStreamingDetails)
			if err != nil {
				return nil, err
			}
			if err := vd.SetScheduledAtUnix(sa); err != nil {
				return nil, err
			}
		case "none":
			vd.SetStatus(api.Archived)
			sa, err := extractScheduledAtUnix(i.LiveStreamingDetails)
			if err != nil {
				return nil, err
			}
			if err := vd.SetScheduledAtUnix(sa); err != nil {
				return nil, err
			}
		case "completed":
			vd.SetStatus(api.Archived)
			sa, err := extractScheduledAtUnix(i.LiveStreamingDetails)
			if err != nil {
				return nil, err
			}
			if err := vd.SetScheduledAtUnix(sa); err != nil {
				return nil, err
			}
		default:
			slog.Error(
				"unexpected liveBroadcastContent",
				"sourceID", i.Id,
				"liveBroadcastContent", i.Snippet.LiveBroadcastContent,
			)
		}

		vds = append(vds, vd)
	}
	return vds, nil
}

func (c *Client) FetchScheduledAtByVideoIDs(ctx context.Context, videoIDs []string) ([]*api.LiveScheduleInfo, error) {
	call := c.ytSvc.Videos.List([]string{"liveStreamingDetails"}).Id(videoIDs...)
	call = call.Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	lsis := make([]*api.LiveScheduleInfo, 0, len(resp.Items))

	for _, i := range resp.Items {
		lsi := api.NewLiveScheduleInfo(i.Id)

		if i.Snippet == nil {
			return nil, fmt.Errorf("snippet is not found")
		}

		// scheduledStartTime
		sa, err := extractScheduledAtUnix(i.LiveStreamingDetails)
		if err != nil {
			return nil, err
		}

		lsi.SetScheduledAtUnix(sa)
		lsis = append(lsis, lsi)
	}

	return lsis, nil
}

func extractChatID(details *youtube.VideoLiveStreamingDetails) string {
	if details == nil {
		return ""
	}

	return details.ActiveLiveChatId
}

func extractScheduledAtUnix(details *youtube.VideoLiveStreamingDetails) (int64, error) {
	if details == nil {
		return 0, nil
	}

	if details.ScheduledStartTime == "" {
		return 0, fmt.Errorf("scheduledStartTime is not found")
	}

	sa, err := synchro.ParseISO[tz.AsiaTokyo](details.ScheduledStartTime)
	if err != nil {
		slog.Error(
			"failed to parse scheduledStartTime",
			"scheduledStartTime", details.ScheduledStartTime,
			slog.Group("extractScheduledAtUnix", "error", err),
		)
		return 0, fmt.Errorf("failed to parse scheduledStartTime: %w", err)
	}

	return sa.Unix(), nil
}
