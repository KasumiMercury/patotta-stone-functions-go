package api

import (
	"context"
	"fmt"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/api/dto"
	repo "github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/youtube"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
	"google.golang.org/api/youtube/v3"
	"log/slog"
)

const (
	PartSnippet              = "snippet"
	PartContentDetails       = "contentDetails"
	PartLiveStreamingDetails = "liveStreamingDetails"
)

const MaxVideoIDs = 50

type YouTubeVideo struct {
	clt repo.Client
}

func NewYouTubeVideo(clt repo.Client) *YouTubeVideo {
	return &YouTubeVideo{clt: clt}
}

func (c *YouTubeVideo) FetchVideoDetailsByVideoIDs(ctx context.Context, videoIDs []string) ([]dto.DetailResponse, error) {
	if len(videoIDs) == 0 {
		return []dto.DetailResponse{}, nil
	}

	idsSlice := make([][]string, 0, len(videoIDs)/MaxVideoIDs+1)
	for i := 0; i < len(videoIDs); i += MaxVideoIDs {
		end := i + MaxVideoIDs
		if end > len(videoIDs) {
			end = len(videoIDs)
		}
		idsSlice = append(idsSlice, videoIDs[i:end])
	}

	vds := make([]dto.DetailResponse, 0, len(videoIDs))

	for _, ids := range idsSlice {
		resp, err := c.clt.VideoList(ctx, []string{PartSnippet, PartContentDetails, PartLiveStreamingDetails}, ids)
		if err != nil {
			return nil, err
		}

		for _, i := range resp.Items {
			vd, err := extractVideoItem(i)
			if err != nil {
				slog.Error(
					"failed to extract video item",
					"sourceID", i.Id,
					"error", err,
				)

				// if any error occurs, skip the item
				// and continue to the next item
				// because the error is not fatal
				continue
			}
			vds = append(vds, *vd)
		}
	}

	return vds, nil
}

func extractVideoItem(i *youtube.Video) (*dto.DetailResponse, error) {
	if i.Snippet == nil {
		return nil, fmt.Errorf("snippet is not found for sourceID: %s", i.Id)
	}

	pa, err := synchro.ParseISO[tz.AsiaTokyo](i.Snippet.PublishedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse publishedAt: %s, %w", i.Snippet.PublishedAt, err)
	}

	sts, cID, sa, err := extractVideoStatus(*i)
	if err != nil {
		return nil, err
	}

	return &dto.DetailResponse{
		Id:          i.Id,
		Title:       i.Snippet.Title,
		Description: i.Snippet.Description,
		Status:      sts,
		PublishedAt: pa,
		ScheduledAt: sa,
		ChatId:      cID,
	}, nil
}

func extractVideoStatus(i youtube.Video) (status.Status, string, synchro.Time[tz.AsiaTokyo], error) {
	switch i.Snippet.LiveBroadcastContent {
	case "live":
		sa, err := extractScheduledAt(i.LiveStreamingDetails)
		if err != nil {
			return status.Live, "", synchro.Time[tz.AsiaTokyo]{}, fmt.Errorf("failed to extract ScheduledAt for live video: %w", err)
		}
		return status.Live, "", sa, nil
	case "upcoming":
		cID := extractChatID(i.LiveStreamingDetails)
		sa, err := extractScheduledAt(i.LiveStreamingDetails)
		if err != nil {
			return status.Upcoming, cID, sa, fmt.Errorf("failed to extract ScheduledAt for upcoming video: %w", err)
		}
		return status.Upcoming, cID, sa, nil
	case "none", "completed":
		sa, err := extractScheduledAt(i.LiveStreamingDetails)
		if err != nil {
			return status.Archived, "", synchro.Time[tz.AsiaTokyo]{}, fmt.Errorf("failed to extract ScheduledAt for archived video: %w", err)
		}
		return status.Archived, "", sa, nil
	case "":
		return status.Undefined, "", synchro.Time[tz.AsiaTokyo]{}, fmt.Errorf("liveBroadcastContent is not found")
	default:
		return status.Undefined, "", synchro.Time[tz.AsiaTokyo]{}, fmt.Errorf("unexpected liveBroadcastContent: %s", i.Snippet.LiveBroadcastContent)
	}
}

func extractChatID(details *youtube.VideoLiveStreamingDetails) string {
	if details == nil {
		return ""
	}

	return details.ActiveLiveChatId
}

func (c *YouTubeVideo) FetchScheduledAtByVideoIDs(ctx context.Context, videoIDs []string) ([]dto.ScheduleResponse, error) {
	if len(videoIDs) == 0 {
		return []dto.ScheduleResponse{}, nil
	}

	idsSlice := make([][]string, 0, len(videoIDs)/MaxVideoIDs+1)
	for i := 0; i < len(videoIDs); i += MaxVideoIDs {
		end := i + MaxVideoIDs
		if end > len(videoIDs) {
			end = len(videoIDs)
		}
		idsSlice = append(idsSlice, videoIDs[i:end])
	}

	vss := make([]dto.ScheduleResponse, 0, len(videoIDs))

	for _, ids := range idsSlice {
		resp, err := c.clt.VideoList(ctx, []string{PartLiveStreamingDetails}, ids)
		if err != nil {
			return nil, err
		}

		for _, i := range resp.Items {
			sa, err := extractScheduledAt(i.LiveStreamingDetails)
			if err != nil {
				slog.Error(
					"failed to extract scheduledStartTime",
					"sourceID", i.Id,
					"error", err,
				)

				// if any error occurs, skip the item
				// and continue to the next item
				// because the error is not fatal
				continue
			}

			vss = append(vss, dto.ScheduleResponse{
				Id:          i.Id,
				ScheduledAt: sa,
			})
		}
	}

	return vss, nil
}

func extractScheduledAt(details *youtube.VideoLiveStreamingDetails) (synchro.Time[tz.AsiaTokyo], error) {
	if details == nil {
		return synchro.Time[tz.AsiaTokyo]{}, nil
	}

	sa, err := synchro.ParseISO[tz.AsiaTokyo](details.ScheduledStartTime)
	if err != nil {
		return synchro.Time[tz.AsiaTokyo]{}, fmt.Errorf("failed to parse scheduledStartTime: %s, %w", details.ScheduledStartTime, err)
	}

	return sa, nil
}
