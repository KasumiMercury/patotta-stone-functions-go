package usecase

import (
	"context"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/model"
	"google.golang.org/api/youtube/v3"
	"log/slog"
)

func FetchChatsByVideoInfo(ctx context.Context, svc *youtube.Service, video model.VideoInfo, l int64) ([]model.YTChat, error) {
	call := svc.LiveChatMessages.List(video.ChatID, []string{"snippet"})

	// If the length is set, set the maximum number of messages to be fetched
	if l > 0 {
		call.MaxResults(l)
	}

	call = call.Context(ctx)

	// Fetch the chats from the YouTube API
	resp, err := call.Do()
	if err != nil {
		slog.Error(
			"Failed to run LiveChatMessages.List",
			slog.Group("fetchChat", "chatId", video.ChatID, "sourceId", video.SourceID, slog.Group("YouTubeAPI", "error", err)),
		)
		return nil, err
	}

	items := resp.Items
	if items == nil || len(items) == 0 {
		slog.Info("No chats returned from the YouTube API",
			slog.Group("fetchChat", "chatId", video.ChatID, "sourceId", video.SourceID),
		)
		return nil, nil
	}

	// Create a slice to store the chats
	chats := make([]model.YTChat, 0, len(resp.Items))
	for _, item := range items {
		pa, err := synchro.ParseISO[tz.AsiaTokyo](item.Snippet.PublishedAt)
		if err != nil {
			slog.Error(
				"Failed to parse the publishedAt",
				slog.Group("fetchChat", "chatId", video.ChatID, "sourceId", video.SourceID,
					slog.Group("formatting", "error", err, "target", item.Snippet.PublishedAt)),
			)
			continue
		}
		chats = append(chats, model.YTChat{
			AuthorChannelID: item.AuthorDetails.ChannelId,
			Message:         item.Snippet.DisplayMessage,
			PublishedAtUnix: pa.Unix(),
			SourceID:        video.SourceID,
		})
	}

	return chats, nil
}
