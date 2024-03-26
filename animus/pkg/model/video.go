package model

import (
	"github.com/uptrace/bun"
	"time"
)

type VideoInfo struct {
	SourceID string `json:"sourceId"`
	ChatID   string `json:"chatId"`
}

type VideoRecord struct {
	bun.BaseModel `bun:"table:videos"`

	SourceID  string    `bun:",type:varchar(255)"`
	Status    string    `bun:",type:varchar(255)"`
	ChatID    string    `bun:",type:varchar(255)"`
	UpdatedAt time.Time `bun:",type:timestamptz"`
}
