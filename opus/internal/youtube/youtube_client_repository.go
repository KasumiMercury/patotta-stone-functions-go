package youtube

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"google.golang.org/api/youtube/v3"
)

type Client interface {
	VideoList(ctx context.Context, part []string, id []string) (*youtube.VideoListResponse, error)
}
