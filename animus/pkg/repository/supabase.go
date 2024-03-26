package repository

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/model"
)

type Supabase interface {
	GetVideoInfoByStatus(ctx context.Context, status []string) ([]model.VideoRecord, error)
	InsertChatRecord(ctx context.Context, record []model.ChatRecord) error
	GetPublishedAtOfLastRecordedChatBySource(ctx context.Context, sourceId string) (int64, error)
	InsertFetchedHistory(ctx context.Context, sourceId string) error
}
