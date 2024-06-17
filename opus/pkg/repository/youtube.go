package repository

import (
	"context"
	"google.golang.org/api/youtube/v3"
)

type YouTube interface {
	FetchVideoDetailsByVideoIDs(ctx context.Context, videoIDs []string) ([]*youtube.Video, error)
	FetchScheduledAtByVideoIDs(ctx context.Context, videoIDs []string) ([]*youtube.Video, error)
}
