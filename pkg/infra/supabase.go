package infra

import (
	"context"
	"database/sql"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log/slog"
)

type SupabaseRepository struct {
	db *bun.DB
}

func NewSupabaseClient(dsn string) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return db, nil
}

func NewSupabaseRepository(dsn string) (*SupabaseRepository, error) {
	db, err := NewSupabaseClient(dsn)
	if err != nil {
		slog.Error(
			"Failed to create a new Supabase client",
			slog.Group("Supabase", "error", err),
		)
		return nil, err
	}

	return &SupabaseRepository{db: db}, nil
}

func (r *SupabaseRepository) GetPublishedAtOfLastRecordedChatBySource(ctx context.Context, sourceId string) (int64, error) {
	records := make([]model.ChatRecord, 0)
	err := r.db.NewSelect().
		Model(&records).
		Where("source_id = ?", sourceId).
		Order("published_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		slog.Error(
			"Failed to get the last recorded chat",
			slog.Group("Supabase", "error", err),
		)
		return 0, err
	}

	if len(records) == 0 {
		return 0, nil
	}

	return records[0].PublishedAt.Unix(), nil
}

func (r *SupabaseRepository) InsertChatRecord(ctx context.Context, record []model.ChatRecord) error {
	_, err := r.db.NewInsert().Model(&record).Exec(ctx)
	if err != nil {
		slog.Error(
			"Failed to insert chat records",
			slog.Group("saveChat", "record", record,
				slog.Group("Supabase", "error", err),
			),
		)
		return err
	}

	return nil
}
