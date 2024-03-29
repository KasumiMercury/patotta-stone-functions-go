package repository

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
)

type Supabase interface {
	GetLastUpdatedUnixOfVideo(ctx context.Context) (int64, error)
	GetRecordsBySourceIDs(ctx context.Context, sourceIDs []string) ([]model.VideoRecord, error)
	InsertVideoRecords(ctx context.Context, records []model.VideoRecord) error
}
