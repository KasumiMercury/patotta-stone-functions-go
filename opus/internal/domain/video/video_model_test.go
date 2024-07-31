package video

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
	"reflect"
	"testing"
	"time"
)

func TestNewVideo(t *testing.T) {
	t.Parallel()

	type args struct {
		channelID   string
		sourceID    string
		title       string
		description string
		chatID      string
		status      status.Status
		publishedAt synchro.Time[tz.AsiaTokyo]
		scheduledAt synchro.Time[tz.AsiaTokyo]
		updatedAt   synchro.Time[tz.AsiaTokyo]
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
				channelID:   "channelID",
				sourceID:    "sourceID",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			want: &Video{
				channelID:   "channelID",
				sourceID:    "sourceID",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "when channelID is empty, return error",
			args: args{
				channelID:   "",
				sourceID:    "sourceID",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when sourceID is empty, return error",
			args: args{
				channelID:   "channelID",
				sourceID:    "",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when title is empty, return error",
			args: args{
				channelID:   "channelID",
				sourceID:    "sourceID",
				title:       "",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when status is undefined, return error",
			args: args{
				channelID:   "channelID",
				sourceID:    "sourceID",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Undefined,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when publishedAtUnix is 0, return error",
			args: args{
				channelID:   "channelID",
				sourceID:    "sourceID",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.Time[tz.AsiaTokyo]{},
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when scheduledAtUnix is less than publishedAtUnix, return error",
			args: args{
				channelID:   "channelID",
				sourceID:    "sourceID",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when scheduledAtUnix is equal to publishedAtUnix, not return error",
			args: args{
				channelID:   "channelID",
				sourceID:    "sourceID",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			want: &Video{
				channelID:   "channelID",
				sourceID:    "sourceID",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
			wantErr: false,
		},
		{
			name: "when updatedAtUnix is 0, not return error",
			args: args{
				channelID:   "channelID",
				sourceID:    "sourceID",
				title:       "title",
				description: "description",
				chatID:      "chatID",
				status:      status.Upcoming,
				publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				updatedAt:   synchro.Time[tz.AsiaTokyo]{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewVideo(tt.args.channelID, tt.args.sourceID, tt.args.title, tt.args.description, tt.args.chatID, tt.args.status, tt.args.publishedAt, tt.args.scheduledAt, tt.args.updatedAt)
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

func TestVideo_PublishedAt(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		publishedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	// Act
	got := v.PublishedAt()

	// Assert
	if got != synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("PublishedAt() got = %v, want %v", got, synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)))
	}
}

func TestVideo_ScheduledAt(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		scheduledAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
	}

	// Act
	got := v.ScheduledAt()

	// Assert
	if got != synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("ScheduledAt() got = %v, want %v", got, synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)))
	}
}

func TestVideo_UpdatedAtUnix(t *testing.T) {
	t.Parallel()
	// Arrange
	v := &Video{
		updatedAt: synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	// Act
	got := v.UpdatedAt()

	// Assert
	if got != synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("UpdatedAtUnix() got = %v, want %v", got, synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)))
	}
}
