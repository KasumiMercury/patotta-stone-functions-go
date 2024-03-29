package infra

import (
	"context"
	"database/sql"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
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

func NewSupabaseRepository(db *bun.DB) *SupabaseRepository {
	return &SupabaseRepository{
		db: db,
	}
}

func (r *SupabaseRepository) GetLastUpdatedUnixOfVideo(ctx context.Context) (int64, error) {
	records := make([]model.VideoRecord, 0)
	err := r.db.NewSelect().
		Model(&records).
		Order("updated_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		slog.Error(
			"Failed to get the last updated video",
			slog.Group("Supabase", "error", err),
		)
		return 0, err
	}

	if len(records) == 0 {
		return 0, nil
	}

	updatedAt := records[0].UpdatedAt
	return updatedAt.Unix(), nil
}
