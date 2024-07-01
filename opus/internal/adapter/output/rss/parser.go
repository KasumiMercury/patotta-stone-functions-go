package rss

//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

import (
	"context"
	"github.com/mmcdole/gofeed"
)

type ParserRepository interface {
	ParseURLWithContext(url string, ctx context.Context) (*gofeed.Feed, error)
}

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
