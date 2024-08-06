package realtime

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/domain/video"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func Test_toDBModel(t *testing.T) {
	t.Parallel()
	type args struct {
		v *video.Video
	}
	tests := []struct {
		name string
		args args
		want *Record
	}{
		{
			name: "normal",
			args: args{
				v: func() *video.Video {
					v, _ := video.NewVideo(
						"channelID",
						"sourceID",
						"title",
						"description",
						"chatID",
						status.Archived,
						synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
						synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
						synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					)
					return v
				}(),
			},
			want: &Record{
				SourceID:    "sourceID",
				Title:       "title",
				Status:      status.Archived.String(),
				ChatID:      "chatID",
				ScheduledAt: timeToPtr(utcToJST(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))),
				UpdatedAt:   utcToJST(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "when scheduledAt is 0",
			args: args{
				v: func() *video.Video {
					v, _ := video.NewVideo(
						"channelID",
						"sourceID",
						"title",
						"description",
						"chatID",
						status.Archived,
						synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
						synchro.In[tz.AsiaTokyo](time.Time{}),
						synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					)
					return v
				}(),
			},
			want: &Record{
				SourceID:    "sourceID",
				Title:       "title",
				Status:      status.Archived.String(),
				ChatID:      "chatID",
				ScheduledAt: nil,
				UpdatedAt:   utcToJST(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if diff := cmp.Diff(tt.want, toDBModel(tt.args.v)); diff != "" {
				t.Errorf("toDBModel() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func timeToPtr(t time.Time) *time.Time {
	return &t
}

func utcToJST(t time.Time) time.Time {
	return t.In(time.FixedZone("JST", 9*60*60))
}
