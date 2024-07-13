package youtube

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func NewYouTubeService(ctx context.Context, apiKey string) (*youtube.Service, error) {
	svc, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return svc, nil
}
