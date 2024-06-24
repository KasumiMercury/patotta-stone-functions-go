package output

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/api"
)

type ApiRepository interface {
	FetchVideoDetailsByVideoIDs(ctx context.Context, videoIDs []string) ([]*api.VideoDetail, error)
}
