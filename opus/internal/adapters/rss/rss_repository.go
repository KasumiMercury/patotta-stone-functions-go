package rss

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapters/rss/dto"
)

type RSSRepository interface {
	FetchRssItems(ctx context.Context, url string, limitUnix int64) ([]dto.Item, error)
}
