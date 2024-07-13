package api

import (
	"context"
	"google.golang.org/api/youtube/v3"
)

type YouTubeImpl struct {
	svc *youtube.Service
}

func NewYouTubeImpl(svc *youtube.Service) *YouTubeImpl {
	return &YouTubeImpl{svc: svc}
}

func (y *YouTubeImpl) VideoList(ctx context.Context, part []string, id string) (*youtube.VideoListResponse, error) {
	call := y.svc.Videos.List(part).Id(id)
	call = call.Context(ctx)

	return call.Do()
}
