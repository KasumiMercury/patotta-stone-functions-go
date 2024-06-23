package output

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/rss"
)

type RSSRepository interface {
	FetchRssItems(ctx context.Context, url string, limitUnix int64) ([]rss.Item, error)
}
