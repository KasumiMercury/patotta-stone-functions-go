package model

import (
	"github.com/uptrace/bun"
	"time"
)

type VideoInfo struct {
	SourceID      string
	Title         string
	UpdatedAtUnix int64
}

type VideoDetail struct {
	SourceID        string
	ChatID          string
	Status          string
	ScheduledAtUnix int64
	PublishedAtUnix int64
}

type VideoRecord struct {
	bun.BaseModel `bun:"table:videos"`

	SourceID    string    `bun:",type:varchar(255)"`
	Title       string    `bun:",type:varchar(255)"`
	Status      string    `bun:",type:varchar(255)"`
	ChatID      string    `bun:",type:varchar(255)"`
	ScheduledAt time.Time `bun:",type:timestamptz"`
	UpdatedAt   time.Time `bun:",type:timestamptz"`
}
