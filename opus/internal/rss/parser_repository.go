package rss

//go:generate mockgen -source=$GOFILE -destination=./../mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"github.com/mmcdole/gofeed"
)

type ParserRepository interface {
	ParseURLWithContext(url string, ctx context.Context) (*gofeed.Feed, error)
}
