package model

import (
	"github.com/uptrace/bun"
	"time"
)

type YTChat struct {
	AuthorChannelID string
	Message         string
	PublishedAtUnix int64
	SourceID        string
}

type ChatRecord struct {
	bun.BaseModel `bun:"table:chats"`

	Message     string    `bun:",pk,type:varchar(255)"`
	IsNegative  bool      `bun:",type:tinyint(1)"`
	SourceID    string    `bun:",type:varchar(255)"`
	PublishedAt time.Time `bun:",type:timestamp"`
}
