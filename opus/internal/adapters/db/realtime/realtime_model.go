package realtime

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/domain/video"
	"github.com/uptrace/bun"
	"time"
)

type Record struct {
	bun.BaseModel `bun:"table:videos"`

	SourceID    string     `bun:",type:varchar(255),unique"`
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
		ScheduledAt: synchroTimeToNillableTime(v.ScheduledAt()),
		UpdatedAt:   v.UpdatedAt().StdTime(),
	}
}

func synchroTimeToNillableTime(t synchro.Time[tz.AsiaTokyo]) *time.Time {
	if t.IsZero() {
		return nil
	}
	tt := t.StdTime()
	return &tt
}
