package output

import (
	"context"
	"github.com/mmcdole/gofeed"
)

type ParserRepository interface {
	ParseURLWithContext(url string, ctx context.Context) (*gofeed.Feed, error)
}
