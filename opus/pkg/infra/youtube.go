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

func (r YouTubeRepository) FetchVideoDetailsByVideoIDs(ctx context.Context, videoIDs []string) ([]*youtube.Video, error) {
	call := r.ytSvc.Videos.List([]string{"snippet", "contentDetails", "liveStreamingDetails"}).Id(videoIDs...)
	call = call.Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (r YouTubeRepository) FetchScheduleAtByVideoIDs(ctx context.Context, videoIDs []string) ([]*youtube.Video, error) {
	call := r.ytSvc.Videos.List([]string{"snippet"}).Id(videoIDs...)
	call = call.Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}
