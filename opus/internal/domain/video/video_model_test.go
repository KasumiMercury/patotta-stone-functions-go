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

func TestVideo_ChannelID(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		channelID: "channelID",
	}

	// Act
	got := v.ChannelID()

	// Assert
	if got != "channelID" {
		t.Errorf("ChannelID() got = %v, want %v", got, "channelID")
	}
}

func TestVideo_SourceID(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		sourceID: "sourceID",
	}

	// Act
	got := v.SourceID()

	// Assert
	if got != "sourceID" {
		t.Errorf("SourceID() got = %v, want %v", got, "sourceID")
	}
}

func TestVideo_Title(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		title: "title",
	}

	// Act
	got := v.Title()

	// Assert
	if got != "title" {
		t.Errorf("Title() got = %v, want %v", got, "title")
	}
}

func TestVideo_Description(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		description: "description",
	}

	// Act
	got := v.Description()

	// Assert
	if got != "description" {
		t.Errorf("Description() got = %v, want %v", got, "description")
	}
}

func TestVideo_ChatID(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		chatID: "chatID",
	}

	// Act
	got := v.ChatID()

	// Assert
	if got != "chatID" {
		t.Errorf("ChatID() got = %v, want %v", got, "chatID")
	}
}

func TestVideo_Status(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		status: status.Upcoming,
	}

	// Act
	got := v.Status()

	// Assert
	if got != status.Upcoming {
		t.Errorf("Status() got = %v, want %v", got, status.Upcoming)
	}
}

func TestVideo_PublishedAtUnix(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}

	// Act
	got := v.PublishedAtUnix()

	// Assert
	if got != time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix() {
		t.Errorf("PublishedAtUnix() got = %v, want %v", got, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
	}
}

func TestVideo_PublishedAt(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		publishedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}

	// Act
	got := v.PublishedAtUTC()

	// Assert
	if got != time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Errorf("PublishedAtUTC() got = %v, want %v", got, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	}
}

func TestVideo_ScheduledAtUnix(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		scheduledAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
	}

	// Act
	got := v.ScheduledAtUnix()

	// Assert
	if got != time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix() {
		t.Errorf("ScheduledAtUnix() got = %v, want %v", got, time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix())
	}
}

func TestVideo_NillableScheduledAt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		v    *Video
		want *time.Time
	}{
		{
			name: "when scheduledAtUnix is 0, return nil",
			v: &Video{
				scheduledAtUnix: 0,
			},
			want: nil,
		},
		{
			name: "when scheduledAtUnix is not 0, return time.Time",
			v: &Video{
				scheduledAtUnix: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Unix(),
			},
			want: timePtr(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.v.NillableScheduledAtUTC()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillableScheduledAtUTC() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideo_UpdatedAtUnix(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		updatedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}

	// Act
	got := v.UpdatedAtUnix()

	// Assert
	if got != time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix() {
		t.Errorf("UpdatedAtUnix() got = %v, want %v", got, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
	}
}

func TestVideo_UpdatedAt(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		updatedAtUnix: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}

	// Act
	got := v.UpdatedAtUTC()

	// Assert
	if got != time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Errorf("UpdatedAtUTC() got = %v, want %v", got, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
