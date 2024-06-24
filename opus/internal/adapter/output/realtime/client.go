package realtime

import (
	"context"
	"database/sql"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/realtime"
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

func (r *Realtime) GetRecordsBySourceIDs(ctx context.Context, sourceIDs []string) ([]realtime.Record, error) {
	records := make([]realtime.Record, 0)
	err := r.db.NewSelect().
		Model(&records).
		Where("source_id IN (?)", bun.In(sourceIDs)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *Realtime) InsertRecords(ctx context.Context, records []realtime.Record) error {
	if _, err := r.db.NewInsert().Model(&records).Exec(ctx); err != nil {
		slog.Error(
			"Failed to insert records into realtime",
			"records", records,
			slog.Group("Realtime", "error", err),
		)
		return err
	}

	return nil
}
