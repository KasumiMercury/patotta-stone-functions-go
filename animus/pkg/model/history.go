package model

import (
	"github.com/uptrace/bun"
	"time"
)

type History struct {
	bun.BaseModel `bun:"table:fetch_chat_history"`

	SourceID  string    `bun:",type:varchar(255)"`
	CreatedAt time.Time `bun:",type:timestamptz"`
}
