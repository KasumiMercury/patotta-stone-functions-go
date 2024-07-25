package realtime

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/domain/video"
	"time"
)

type RealtimeRepository interface {
	UpsertRecords(ctx context.Context, videos []video.Video) error
	GetRecordsBySourceIDs(ctx context.Context, sourceIDs []string) ([]*Record, error)
	InsertRecords(ctx context.Context, videos []video.Video) error
	UpdateRecords(ctx context.Context, videos []video.Video) error
	GetLastUpdatedUnixOfVideo(ctx context.Context) (int64, error)
	UpdateScheduledAtBySourceID(ctx context.Context, sourceID string, scheduledAt time.Time) error
	UpdateStatusBySourceID(ctx context.Context, sourceID string, status string) error
}
