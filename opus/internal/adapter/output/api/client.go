package api

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	svc *youtube.Service
}

func NewYouTubeClient(c config.Config) (*Client, error) {
	svc, err := youtube.NewService(context.Background(), option.WithAPIKey(c.Api.ApiKey))
	if err != nil {
		return nil, err
	}

	return &Client{svc: svc}, nil
}

func (y *Client) VideoList(ctx context.Context, part []string, id []string) (*youtube.VideoListResponse, error) {
	call := y.svc.Videos.List(part).Id(id...)
	call = call.Context(ctx)

	return call.Do()
}
