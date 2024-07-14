package video

import (
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/status"
	"reflect"
	"testing"
	"time"
)

func TestVideo_NillableChatID(t *testing.T) {
	type fields struct {
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

	t.Parallel()

	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{
			name: "chatID is empty",
			fields: fields{
				chatID: "",
			},
			want: nil,
		},
		{
			name: "chatID is not empty",
			fields: fields{
				chatID: "chatID",
			},
			want: strPtr("chatID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{
				channelID:       tt.fields.channelID,
				sourceID:        tt.fields.sourceID,
				title:           tt.fields.title,
				description:     tt.fields.description,
				chatID:          tt.fields.chatID,
				status:          tt.fields.status,
				publishedAtUnix: tt.fields.publishedAtUnix,
				scheduledAtUnix: tt.fields.scheduledAtUnix,
				updatedAtUnix:   tt.fields.updatedAtUnix,
			}
			if got := v.NillableChatID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillableChatID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideo_NillablePublishedAt(t *testing.T) {
	type fields struct {
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
		name   string
		fields fields
		want   *time.Time
	}{
		{
			name: "publishedAtUnix is 0",
			fields: fields{
				publishedAtUnix: 0,
			},
			want: nil,
		},
		{
			name: "publishedAtUnix is not 0",
			fields: fields{
				publishedAtUnix: 1610000000,
			},
			want: timePtr(time.Unix(1610000000, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{
				channelID:       tt.fields.channelID,
				sourceID:        tt.fields.sourceID,
				title:           tt.fields.title,
				description:     tt.fields.description,
				chatID:          tt.fields.chatID,
				status:          tt.fields.status,
				publishedAtUnix: tt.fields.publishedAtUnix,
				scheduledAtUnix: tt.fields.scheduledAtUnix,
				updatedAtUnix:   tt.fields.updatedAtUnix,
			}
			if got := v.NillablePublishedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillablePublishedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideo_NillableScheduledAt(t *testing.T) {
	type fields struct {
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
		name   string
		fields fields
		want   *time.Time
	}{
		{
			name: "scheduledAtUnix is 0",
			fields: fields{
				scheduledAtUnix: 0,
			},
			want: nil,
		},
		{
			name: "scheduledAtUnix is not 0",
			fields: fields{
				scheduledAtUnix: 1610000000,
			},
			want: timePtr(time.Unix(1610000000, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{
				channelID:       tt.fields.channelID,
				sourceID:        tt.fields.sourceID,
				title:           tt.fields.title,
				description:     tt.fields.description,
				chatID:          tt.fields.chatID,
				status:          tt.fields.status,
				publishedAtUnix: tt.fields.publishedAtUnix,
				scheduledAtUnix: tt.fields.scheduledAtUnix,
				updatedAtUnix:   tt.fields.updatedAtUnix,
			}
			if got := v.NillableScheduledAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillableScheduledAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideo_NillableUpdatedAt(t *testing.T) {
	type fields struct {
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
		name   string
		fields fields
		want   *time.Time
	}{
		{
			name: "updatedAtUnix is 0",
			fields: fields{
				updatedAtUnix: 0,
			},
			want: nil,
		},
		{
			name: "updatedAtUnix is not 0",
			fields: fields{
				updatedAtUnix: 1610000000,
			},
			want: timePtr(time.Unix(1610000000, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{
				channelID:       tt.fields.channelID,
				sourceID:        tt.fields.sourceID,
				title:           tt.fields.title,
				description:     tt.fields.description,
				chatID:          tt.fields.chatID,
				status:          tt.fields.status,
				publishedAtUnix: tt.fields.publishedAtUnix,
				scheduledAtUnix: tt.fields.scheduledAtUnix,
				updatedAtUnix:   tt.fields.updatedAtUnix,
			}
			if got := v.NillableUpdatedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillableUpdatedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
