package repository

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
)

type Rss interface {
	FetchUpdatedRssItems(ctx context.Context, url string, threshold int64) ([]model.Rss, error)
}
