package output

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapter/output/api/dto"
)

type ApiRepository interface {
	FetchVideoDetailsByVideoIDs(ctx context.Context, videoIDs []string) ([]dto.DetailResponse, error)
	FetchScheduledAtByVideoIDs(ctx context.Context, videoIDs []string) ([]dto.ScheduleResponse, error)
}
