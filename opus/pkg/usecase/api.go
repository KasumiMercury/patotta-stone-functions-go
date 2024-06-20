package usecase

import (
	"context"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/repository"
	"log/slog"
)

type Api interface {
	VideoDetailsMap(ctx context.Context, sourceIDs []string) (map[string]model.VideoDetail, error)
	VideoScheduledAtUnixMap(ctx context.Context, sourceIDs []string) (map[string]int64, error)
}

type apiUsecase struct {
	ytRepo repository.YouTube
}

func NewApiUsecase(ytRepo repository.YouTube) Api {
	return &apiUsecase{
		ytRepo: ytRepo,
	}
}

func (u *apiUsecase) VideoDetailsMap(ctx context.Context, sourceIDs []string) (map[string]model.VideoDetail, error) {
	// fetch video details from YouTube api
	items, err := u.ytRepo.FetchVideoDetailsByVideoIDs(ctx, sourceIDs)
	if err != nil {
		slog.Error(
			"failed to fetch video details",
			"sourceIDs", sourceIDs,
			slog.Group("YouTubeAPI", "error", err),
		)
		return nil, err
	}

	// format video details to map
	dtlMap := make(map[string]model.VideoDetail, len(items))
	for _, i := range items {
		var dtl model.VideoDetail
		dtl.SourceID = i.Id

		pa, err := synchro.ParseISO[tz.AsiaTokyo](i.Snippet.PublishedAt)
		if err != nil {
			slog.Error(
				"failed to parse publishedAt",
				"sourceID", i.Id,
				"publishedAt", i.Snippet.PublishedAt,
				slog.Group("formatVideoDetails", "error", err),
			)
			return nil, err
		}
		dtl.PublishedAtUnix = pa.Unix()

		switch i.Snippet.LiveBroadcastContent {
		case "live":
			dtl.Status = "live"
			dtl.ChatID = i.LiveStreamingDetails.ActiveLiveChatId
			sa, err := synchro.ParseISO[tz.AsiaTokyo](i.LiveStreamingDetails.ScheduledStartTime)
			if err != nil {
				slog.Error(
					"failed to parse actualStartTime",
					"sourceID", i.Id,
					"actualStartTime", i.LiveStreamingDetails.ActualStartTime,
					slog.Group("formatVideoDetails", "error", err),
				)
				continue
			}
			dtl.ScheduledAtUnix = sa.Unix()
		case "upcoming":
			dtl.Status = "upcoming"
			dtl.ChatID = i.LiveStreamingDetails.ActiveLiveChatId
			sa, err := synchro.ParseISO[tz.AsiaTokyo](i.LiveStreamingDetails.ScheduledStartTime)
			if err != nil {
				slog.Error(
					"failed to parse actualStartTime",
					"sourceID", i.Id,
					"actualStartTime", i.LiveStreamingDetails.ActualStartTime,
					slog.Group("formatVideoDetails", "error", err),
				)
				continue
			}
			dtl.ScheduledAtUnix = sa.Unix()
		case "none":
			dtl.Status = "archived"
			if i.LiveStreamingDetails != nil {
				sa, err := synchro.ParseISO[tz.AsiaTokyo](i.LiveStreamingDetails.ScheduledStartTime)
				if err != nil {
					slog.Error(
						"failed to parse actualStartTime",
						"sourceID", i.Id,
						"actualStartTime", i.LiveStreamingDetails.ActualStartTime,
						slog.Group("formatVideoDetails", "error", err),
					)
					continue
				}
				dtl.ScheduledAtUnix = sa.Unix()
			}
		default:
			slog.Error(
				"unexpected liveBroadcastContent",
				"sourceID", i.Id,
				"liveBroadcastContent", i.Snippet.LiveBroadcastContent,
			)
			continue
		}

		dtlMap[i.Id] = dtl
	}

	return dtlMap, nil
}

func (u *apiUsecase) VideoScheduledAtUnixMap(ctx context.Context, sourceIDs []string) (map[string]int64, error) {
	// fetch scheduled time of upcoming videos
	items, err := u.ytRepo.FetchScheduledAtByVideoIDs(ctx, sourceIDs)
	if err != nil {
		slog.Error(
			"failed to fetch scheduled time",
			"sourceIDs", sourceIDs,
			slog.Group("YouTubeAPI", "error", err),
		)
		return nil, err
	}

	// format scheduled time to map
	saMap := make(map[string]int64, len(items))
	for _, i := range items {
		sa, err := synchro.ParseISO[tz.AsiaTokyo](i.LiveStreamingDetails.ScheduledStartTime)
		if err != nil {
			slog.Error(
				"failed to parse scheduledStartTime",
				"sourceID", i.Id,
				"scheduledStartTime", i.LiveStreamingDetails.ScheduledStartTime,
				slog.Group("formatScheduledTime", "error", err),
			)
			continue
		}
		saMap[i.Id] = sa.Unix()
	}

	return saMap, nil
}
