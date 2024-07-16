package api

import (
	"context"
	"fmt"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/api"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/status"
	"google.golang.org/api/youtube/v3"
	"log/slog"
)

type YouTubeVideo struct {
	clt output.Client
}

func NewYouTubeVideo(clt output.Client) (*YouTubeVideo, error) {
	return &YouTubeVideo{clt: clt}, nil
}

func (c *YouTubeVideo) FetchVideoDetailsByVideoIDs(ctx context.Context, videoIDs []string) ([]api.VideoDetail, error) {
	resp, err := c.clt.VideoList(ctx, []string{"snippet", "contentDetails", "liveStreamingDetails"}, videoIDs)
	if err != nil {
		return nil, err
	}

	vds := make([]api.VideoDetail, 0, len(resp.Items))

	for _, i := range resp.Items {
		vd, err := c.extractVideoItem(ctx, i)
		if err != nil {
			slog.Error(
				"failed to extract video item",
				"sourceID", i.Id,
				"error", err,
			)
			continue
		}
		vds = append(vds, *vd)
	}
	return vds, nil
}

func (c *YouTubeVideo) extractVideoItem(ctx context.Context, i *youtube.Video) (*api.VideoDetail, error) {
	if i.Snippet == nil {
		return nil, fmt.Errorf("snippet is not found for sourceID: %s", i.Id)
	}

	paTime, err := synchro.ParseISO[tz.AsiaTokyo](i.Snippet.PublishedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'publishedAt' for video ID: %s: %w", i.Id, err)
	}

	sts, cID, sa, err := c.extractVideoStatus(*i)
	if err != nil {
		return nil, fmt.Errorf("failed to extract video status for sourceID: %s: %w", i.Id, err)
	}

	return api.NewVideoDetail(i.Id, cID, sts, paTime.Unix(), sa)
}

func (c *YouTubeVideo) extractVideoStatus(i youtube.Video) (status.Status, string, int64, error) {
	switch i.Snippet.LiveBroadcastContent {
	case "live":
		sa, err := extractScheduledAtUnix(i.LiveStreamingDetails)
		return status.Live, "", sa, err
	case "upcoming":
		cID := extractChatID(i.LiveStreamingDetails)
		sa, err := extractScheduledAtUnix(i.LiveStreamingDetails)
		return status.Upcoming, cID, sa, err
	case "none", "completed":
		sa, err := extractScheduledAtUnix(i.LiveStreamingDetails)
		return status.Archived, "", sa, err
	default:
		return status.Undefined, "", 0, fmt.Errorf("unexpected liveBroadcastContent: %s", i.Snippet.LiveBroadcastContent)
	}
}

func (c *YouTubeVideo) FetchScheduledAtByVideoIDs(ctx context.Context, videoIDs []string) ([]api.LiveScheduleInfo, error) {
	resp, err := c.clt.VideoList(ctx, []string{"liveStreamingDetails"}, videoIDs)
	if err != nil {
		return nil, err
	}

	lsis := make([]api.LiveScheduleInfo, 0, len(resp.Items))

	for _, i := range resp.Items {
		lsi := api.NewLiveScheduleInfo(i.Id)

		// scheduledStartTime
		sa, err := extractScheduledAtUnix(i.LiveStreamingDetails)
		if err != nil {
			return nil, err
		}

		lsi.SetScheduledAtUnix(sa)
		lsis = append(lsis, *lsi)
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
