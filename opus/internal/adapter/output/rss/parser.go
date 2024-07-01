package rss

import (
	"context"
	"github.com/mmcdole/gofeed"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseURLWithContext(url string, ctx context.Context) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()

	feed, err := parser.ParseURLWithContext(url, ctx)
	return feed, err
}
