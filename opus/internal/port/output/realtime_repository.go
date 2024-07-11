package output

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapter/output/db/realtime"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/video"
	"time"
)

type RealtimeRepository interface {
	UpsertRecords(ctx context.Context, videos []video.Video) error
	GetRecordsBySourceIDs(ctx context.Context, sourceIDs []string) ([]*realtime.Record, error)
	InsertRecords(ctx context.Context, videos []video.Video) error
	UpdateRecords(ctx context.Context, videos []video.Video) error
	GetLastUpdatedUnixOfVideo(ctx context.Context) (int64, error)
	UpdateScheduledAtBySourceID(ctx context.Context, sourceID string, scheduledAt time.Time) error
	UpdateStatusBySourceID(ctx context.Context, sourceID string, status string) error
}
