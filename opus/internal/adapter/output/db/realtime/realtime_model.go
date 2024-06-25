package realtime

import (
	"github.com/uptrace/bun"
	"time"
)

type Record struct {
	bun.BaseModel `bun:"table:videos"`

	SourceID    string    `bun:",type:varchar(255)"`
	Title       string    `bun:",type:varchar(255)"`
	Status      string    `bun:",type:varchar(255)"`
	ChatID      string    `bun:",type:varchar(255)"`
	ScheduledAt time.Time `bun:",type:timestamptz"`
	UpdatedAt   time.Time `bun:",type:timestamptz"`
}
