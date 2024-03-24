package repository

import (
	"context"
	"google.golang.org/api/youtube/v3"
)

type YouTube interface {
	FetchChatsByChatID(ctx context.Context, chatID string, l int64) (*youtube.LiveChatMessageListResponse, error)
}
