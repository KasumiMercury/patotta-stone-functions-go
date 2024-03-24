package service

import (
	"context"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/repository"
	"log/slog"
)

type Chat interface {
	FetchChatsByVideoInfo(ctx context.Context, videoInfo model.VideoInfo, l int64) ([]model.YTChat, error)
}

type chatService struct {
	ytRepo repository.YouTube
}

func NewChatService(ytRepo repository.YouTube) Chat {
	return &chatService{
		ytRepo: ytRepo,
	}
}

func (s *chatService) FetchChatsByVideoInfo(ctx context.Context, videoInfo model.VideoInfo, l int64) ([]model.YTChat, error) {
	// Fetch chats from the static target video
	resp, err := s.ytRepo.FetchChatsByChatID(ctx, videoInfo.ChatID, l)
	if err != nil {
		slog.Error(
			"Failed to fetch chats",
			slog.Group("fetchChat", "chatId", videoInfo.ChatID, "sourceId", videoInfo.SourceID,
				slog.Group("YouTubeAPI", "error", err),
			),
		)
		return nil, err
	}

	items := resp.Items
	if items == nil || len(items) == 0 {
		return nil, nil
	}

	// Create a slice to store the chats
	chats := make([]model.YTChat, 0, len(resp.Items))
	for _, item := range items {
		pa, err := synchro.ParseISO[tz.AsiaTokyo](item.Snippet.PublishedAt)
		if err != nil {
			slog.Error(
				"Failed to parse the publishedAt",
				slog.Group("fetchChat", "chatId", videoInfo.ChatID, "sourceId", videoInfo.SourceID,
					slog.Group("formatting", "error", err, "target", item.Snippet.PublishedAt)),
			)
			continue
		}
		chats = append(chats, model.YTChat{
			AuthorChannelID: item.Snippet.AuthorChannelId,
			Message:         item.Snippet.DisplayMessage,
			PublishedAtUnix: pa.Unix(),
			SourceID:        videoInfo.SourceID,
		})
	}

	return chats, nil
}
