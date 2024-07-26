package video

import (
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
	"reflect"
	"testing"
	"time"
)

func TestNewVideo(t *testing.T) {
	t.Parallel()

	type args struct {
		channelID       string
		sourceID        string
		title           string
		description     string
		chatID          string
		status          status.Status
		publishedAtUnix int64
		scheduledAtUnix int64
		updatedAtUnix   int64
	}
	tests := []struct {
		name    string
		args    args
		want    *Video
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				channelID:       "channelID",
				sourceID:        "sourceID",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want: &Video{
				channelID:       "channelID",
				sourceID:        "sourceID",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			wantErr: false,
		},
		{
			name: "when channelID is empty, return error",
			args: args{
				channelID:       "",
				sourceID:        "sourceID",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when sourceID is empty, return error",
			args: args{
				channelID:       "channelID",
				sourceID:        "",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when title is empty, return error",
			args: args{
				channelID:       "channelID",
				sourceID:        "sourceID",
				title:           "",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when status is undefined, return error",
			args: args{
				channelID:       "channelID",
				sourceID:        "sourceID",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Undefined,
				publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when publishedAtUnix is 0, return error",
			args: args{
				channelID:       "channelID",
				sourceID:        "sourceID",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: 0,
				scheduledAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when scheduledAtUnix is less than publishedAtUnix, return error",
			args: args{
				channelID:       "channelID",
				sourceID:        "sourceID",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when scheduledAtUnix is equal to publishedAtUnix, not return error",
			args: args{
				channelID:       "channelID",
				sourceID:        "sourceID",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want: &Video{
				channelID:       "channelID",
				sourceID:        "sourceID",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			},
			wantErr: false,
		},
		{
			name: "when updatedAtUnix is 0, not return error",
			args: args{
				channelID:       "channelID",
				sourceID:        "sourceID",
				title:           "title",
				description:     "description",
				chatID:          "chatID",
				status:          status.Upcoming,
				publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				scheduledAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				updatedAtUnix:   0,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewVideo(tt.args.channelID, tt.args.sourceID, tt.args.title, tt.args.description, tt.args.chatID, tt.args.status, tt.args.publishedAtUnix, tt.args.scheduledAtUnix, tt.args.updatedAtUnix)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVideo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
