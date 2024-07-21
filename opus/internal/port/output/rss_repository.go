package output

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapter/output/rss/dto"
)

type RSSRepository interface {
	FetchRssItems(ctx context.Context, url string, limitUnix int64) ([]dto.Item, error)
}
