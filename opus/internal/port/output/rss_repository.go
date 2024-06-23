package output

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/rss"
)

type RSSRepository interface {
	GetRssItems(ctx context.Context, url string) ([]rss.Item, error)
}
