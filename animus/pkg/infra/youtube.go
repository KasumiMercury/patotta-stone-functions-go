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

func NewYouTubeRepository(svc *youtube.Service) *YouTubeRepository {
	return &YouTubeRepository{ytSvc: svc}
}

func (r *YouTubeRepository) FetchChatsByChatID(ctx context.Context, chatID string, maxResults int64) (*youtube.LiveChatMessageListResponse, bool, error) {
	call := r.ytSvc.LiveChatMessages.List(chatID, []string{"snippet"})

	// If the length is set, set the maximum number of messages to be fetched
	if maxResults > 0 {
		call.MaxResults(maxResults)
	}

	call = call.Context(ctx)

	// Fetch the chats from the YouTube API
	resp, err := call.Do()
	if err != nil {
		// If video is archived and an error occurs, recover the error and notify archiving
		if err.Error() == "googleapi: Error 403: The live chat is no longer live., liveChatEnded" {
			return nil, true, nil
		}
		// Otherwise, log the error and return the error as common error
		slog.Error(
			"Failed to run LiveChatMessages.List",
			slog.Group("fetchChat", "chatId", chatID, slog.Group("YouTubeAPI", "error", err)),
		)
		return nil, false, err
	}
	return resp, false, nil
}
