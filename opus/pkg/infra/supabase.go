package infra

import (
	"context"
	"database/sql"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log/slog"
	"time"
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

func (r *SupabaseRepository) GetRecordsBySourceIDs(ctx context.Context, sourceIDs []string) ([]model.VideoRecord, error) {
	records := make([]model.VideoRecord, 0)
	err := r.db.NewSelect().
		Model(&records).
		Where("source_id IN (?)", bun.In(sourceIDs)).
		Scan(ctx)
	if err != nil {
		slog.Error(
			"Failed to get video records by source IDs",
			"sourceIDs", sourceIDs,
			slog.Group("Supabase", "error", err),
		)
		return nil, err
	}

	return records, nil
}

func (r *SupabaseRepository) InsertVideoRecords(ctx context.Context, records []model.VideoRecord) error {
	_, err := r.db.NewInsert().Model(&records).Exec(ctx)
	if err != nil {
		slog.Error(
			"Failed to insert video records",
			slog.Group("Supabase", "error", err),
		)
		return err
	}

	return nil
}

func (r *SupabaseRepository) UpdateScheduledAtBySourceID(ctx context.Context, sourceID string, scheduledAt time.Time) error {
	_, err := r.db.NewUpdate().
		Model(&model.VideoRecord{ScheduledAt: scheduledAt}).
		Where("source_id = ?", sourceID).
		Exec(ctx)
	if err != nil {
		slog.Error(
			"Failed to update scheduledAt",
			"sourceID", sourceID,
			"scheduledAt", scheduledAt,
			slog.Group("Supabase", "error", err),
		)
		return err
	}

	return nil
}

func (r *SupabaseRepository) UpdateStatusBySourceID(ctx context.Context, sourceID string, status string) error {
	_, err := r.db.NewUpdate().
		Model(&model.VideoRecord{Status: status}).
		Where("source_id = ?", sourceID).
		Exec(ctx)
	if err != nil {
		slog.Error(
			"Failed to update status",
			"sourceID", sourceID,
			"status", status,
			slog.Group("Supabase", "error", err),
		)
		return err
	}

	return nil
}
