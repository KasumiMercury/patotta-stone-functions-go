package realtime

import (
	"context"
	"database/sql"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/domain/video"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log/slog"
)

type Realtime struct {
	db *bun.DB
}

func NewRealtimeClient(dsn string) (*Realtime, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	return &Realtime{db: db}, nil
}

func (r *Realtime) UpsertRecords(ctx context.Context, videos []video.Video) error {
	rec := make([]*Record, 0, len(videos))
	for _, v := range videos {
		rec = append(rec, toDBModel(&v))
	}

	if _, err := r.db.NewInsert().Model(&rec).
		On("conflict (source_id) do update").
		Set("source_id = EXCLUDED.source_id").
		Exec(ctx); err != nil {
		slog.Error(
			"Failed to upsert records into realtime",
			"videos", videos,
			slog.Group("Realtime", "error", err),
		)
		return err
	}

	return nil
}

func (r *Realtime) GetRecordsBySourceIDs(ctx context.Context, sourceIDs []string) ([]*Record, error) {
	records := make([]*Record, 0)
	err := r.db.NewSelect().
		Model(&records).
		Where("source_id IN (?)", bun.In(sourceIDs)).
		Scan(ctx)
	if err != nil {
		slog.Error(
			"Failed to get records by source IDs",
			"sourceIDs", sourceIDs,
			slog.Group("Realtime", "error", err),
		)
		return nil, err
	}

	return records, nil
}

func (r *Realtime) GetLastUpdatedUnixOfVideo(ctx context.Context) (int64, error) {
	records := make([]Record, 0)
	err := r.db.NewSelect().
		Model(&records).
		Order("updated_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		slog.Error(
			"Failed to get the last updated video",
			slog.Group("Realtime", "error", err),
		)
		return 0, err
	}

	if len(records) == 0 {
		return 0, nil
	}

	updatedAt := records[0].UpdatedAt
	return updatedAt.Unix(), nil
}
