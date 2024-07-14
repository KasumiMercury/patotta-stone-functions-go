package api

import (
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/status"
	"reflect"
	"testing"
	"time"
)

func TestLiveScheduleInfo_NillableScheduledAt(t *testing.T) {
	type fields struct {
		sourceID        string
		scheduledAtUnix int64
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
			lsi := &LiveScheduleInfo{
				sourceID:        tt.fields.sourceID,
				scheduledAtUnix: tt.fields.scheduledAtUnix,
			}
			if got := lsi.NillableScheduledAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillableScheduledAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideoDetail_NillableChatID(t *testing.T) {
	type fields struct {
		sourceID        string
		chatID          string
		status          status.Status
		publishedAtUnix int64
		scheduledAtUnix int64
	}
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
			vd := &VideoDetail{
				sourceID:        tt.fields.sourceID,
				chatID:          tt.fields.chatID,
				status:          tt.fields.status,
				publishedAtUnix: tt.fields.publishedAtUnix,
				scheduledAtUnix: tt.fields.scheduledAtUnix,
			}
			if got := vd.NillableChatID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillableChatID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideoDetail_NillablePublishedAt(t *testing.T) {
	type fields struct {
		sourceID        string
		chatID          string
		status          status.Status
		publishedAtUnix int64
		scheduledAtUnix int64
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
			vd := &VideoDetail{
				sourceID:        tt.fields.sourceID,
				chatID:          tt.fields.chatID,
				status:          tt.fields.status,
				publishedAtUnix: tt.fields.publishedAtUnix,
				scheduledAtUnix: tt.fields.scheduledAtUnix,
			}
			if got := vd.NillablePublishedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillablePublishedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideoDetail_NillableScheduledAt(t *testing.T) {
	type fields struct {
		sourceID        string
		chatID          string
		status          status.Status
		publishedAtUnix int64
		scheduledAtUnix int64
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
			vd := &VideoDetail{
				sourceID:        tt.fields.sourceID,
				chatID:          tt.fields.chatID,
				status:          tt.fields.status,
				publishedAtUnix: tt.fields.publishedAtUnix,
				scheduledAtUnix: tt.fields.scheduledAtUnix,
			}
			if got := vd.NillableScheduledAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillableScheduledAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideoDetail_SetScheduledAtUnix(t *testing.T) {
	type fields struct {
		sourceID        string
		chatID          string
		status          status.Status
		publishedAtUnix int64
		scheduledAtUnix int64
	}
	type args struct {
		scheduledAtUnix int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "publishedAtUnix is already set",
			fields: fields{
				publishedAtUnix: 1610000000,
			},
			args: args{
				scheduledAtUnix: 1610000000,
			},
			wantErr: false,
		},
		{
			name:   "publishedAtUnix is not set",
			fields: fields{},
			args: args{
				scheduledAtUnix: 1610000000,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vd := &VideoDetail{
				sourceID:        tt.fields.sourceID,
				chatID:          tt.fields.chatID,
				status:          tt.fields.status,
				publishedAtUnix: tt.fields.publishedAtUnix,
				scheduledAtUnix: tt.fields.scheduledAtUnix,
			}
			if err := vd.SetScheduledAtUnix(tt.args.scheduledAtUnix); (err != nil) != tt.wantErr {
				t.Errorf("SetScheduledAtUnix() error = %v, wantErr %v", err, tt.wantErr)
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
