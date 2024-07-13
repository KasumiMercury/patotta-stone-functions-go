package api

import (
	"context"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	svc *youtube.Service
}

func NewYouTubeClient(svc *youtube.Service) *Client {
	return &Client{svc: svc}
}

func (y *Client) VideoList(ctx context.Context, part []string, id []string) (*youtube.VideoListResponse, error) {
	call := y.svc.Videos.List(part).Id(id...)
	call = call.Context(ctx)

	return call.Do()
}
