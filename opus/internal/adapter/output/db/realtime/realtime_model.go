package realtime

import (
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/video"
	"github.com/uptrace/bun"
	"time"
)

type Record struct {
	bun.BaseModel `bun:"table:videos"`

	SourceID    string     `bun:",type:varchar(255)"`
	Title       string     `bun:",type:varchar(255)"`
	Status      string     `bun:",type:varchar(255)"`
	ChatID      string     `bun:",type:varchar(255)"`
	ScheduledAt *time.Time `bun:",type:timestamptz"`
	UpdatedAt   *time.Time `bun:",type:timestamptz"`
}

func toDBModel(v *video.Video) *Record {
	return &Record{
		SourceID:    v.SourceID(),
		Title:       v.Title(),
		Status:      v.Status(),
		ChatID:      v.ChatID(),
		ScheduledAt: v.NillableScheduledAt(),
		UpdatedAt:   v.NillableUpdatedAt(),
	}
}
