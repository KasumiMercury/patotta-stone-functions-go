package infra

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log/slog"
)

type YouTubeRepository struct {
	ytSvc *youtube.Service
}

func NewYouTubeService(ctx context.Context, apiKey string) (*youtube.Service, error) {
	return youtube.NewService(ctx, option.WithAPIKey(apiKey))
}

func NewYouTubeRepository(ctx context.Context, apiKey string) (*YouTubeRepository, error) {
	ytSvc, err := NewYouTubeService(ctx, apiKey)
	if err != nil {
		slog.Error(
			"Failed to create a new YouTube service",
			slog.Group("YouTubeAPI", "error", err),
		)
		return nil, err
	}

	return &YouTubeRepository{ytSvc: ytSvc}, nil
}

func (r *YouTubeRepository) FetchChatsByChatID(ctx context.Context, chatID string, l int64) (*youtube.LiveChatMessageListResponse, error) {
	call := r.ytSvc.LiveChatMessages.List(chatID, []string{"snippet"})

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
			slog.Group("fetchChat", "chatId", chatID, slog.Group("YouTubeAPI", "error", err)),
		)
		return nil, err
	}
	return resp, nil
}
