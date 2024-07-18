package output

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapter/output/api"
	apiDomain "github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/api"
)

type ApiRepository interface {
	FetchVideoDetailsByVideoIDs(ctx context.Context, videoIDs []string) ([]api.DetailResponse, error)
	FetchScheduledAtByVideoIDs(ctx context.Context, videoIDs []string) ([]apiDomain.LiveScheduleInfo, error)
}
