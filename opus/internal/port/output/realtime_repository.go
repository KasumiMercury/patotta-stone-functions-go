package output

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapter/output/db/realtime"
	"time"
)

type RealtimeRepository interface {
	GetRecordsBySourceIDs(ctx context.Context, sourceIDs []string) ([]*realtime.Record, error)
	InsertRecords(ctx context.Context, records []*realtime.Record) error
	UpdateRecords(ctx context.Context, records []*realtime.Record) error
	GetLastUpdatedUnixOfVideo(ctx context.Context) (int64, error)
	UpdateScheduledAtBySourceID(ctx context.Context, sourceID string, scheduledAt time.Time) error
	UpdateStatusBySourceID(ctx context.Context, sourceID string, status string) error
}
