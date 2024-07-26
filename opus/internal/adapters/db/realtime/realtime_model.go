package realtime

import (
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/domain/video"
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
	UpdatedAt   time.Time  `bun:",type:timestamptz"`
}

func toDBModel(v *video.Video) *Record {
	return &Record{
		SourceID:    v.SourceID(),
		Title:       v.Title(),
		Status:      v.Status().String(),
		ChatID:      v.ChatID(),
		ScheduledAt: unixToNillableTime(v.ScheduledAtUnix()),
		UpdatedAt:   *unixToNillableTime(v.UpdatedAtUnix()),
	}
}

func unixToNillableTime(unix int64) *time.Time {
	if unix == 0 {
		return nil
	}
	t := time.Unix(unix, 0).UTC()
	return &t
}
