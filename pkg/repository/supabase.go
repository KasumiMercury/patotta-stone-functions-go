package repository

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/model"
)

type Supabase interface {
	InsertChatRecord(ctx context.Context, record []model.ChatRecord) error
	GetPublishedAtOfLastRecordedChatBySource(ctx context.Context, sourceId string) (int64, error)
}
