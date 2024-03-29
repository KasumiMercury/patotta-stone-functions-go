package infra

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
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
