package model

import (
	"github.com/uptrace/bun"
	"time"
)

type VideoRecord struct {
	bun.BaseModel `bun:"table:videos"`

	SourceID  string    `bun:",type:varchar(255)"`
	Title     string    `bun:",type:varchar(255)"`
	Status    string    `bun:",type:varchar(255)"`
	ChatID    string    `bun:",type:varchar(255)"`
	UpdatedAt time.Time `bun:",type:timestamptz"`
}
