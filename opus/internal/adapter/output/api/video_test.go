package api

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/youtube/v3"
	"reflect"
	"testing"
)

func TestYouTubeVideo_FetchScheduledAtByVideoIDs(t *testing.T) {
	type fields struct {
		clt *Client
	}
	type args struct {
		ctx      context.Context
		videoIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []api.LiveScheduleInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &YouTubeVideo{
				clt: tt.fields.clt,
			}
			got, err := c.FetchScheduledAtByVideoIDs(tt.args.ctx, tt.args.videoIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchScheduledAtByVideoIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchScheduledAtByVideoIDs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYouTubeVideo_FetchVideoDetailsByVideoIDs(t *testing.T) {
	type fields struct {
		clt *Client
	}
	type args struct {
		ctx      context.Context
		videoIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []api.VideoDetail
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &YouTubeVideo{
				clt: tt.fields.clt,
			}
			got, err := c.FetchVideoDetailsByVideoIDs(tt.args.ctx, tt.args.videoIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchVideoDetailsByVideoIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchVideoDetailsByVideoIDs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractScheduledAtUnix(t *testing.T) {
	type args struct {
		details *youtube.VideoLiveStreamingDetails
	}

	t.Parallel()

	tests := map[string]struct {
		args    args
		want    int64
		wantErr bool
	}{
		"details is nil": {
			args:    args{details: nil},
			want:    0,
			wantErr: false,
		},
		"details is not nil": {
			args: args{details: &youtube.VideoLiveStreamingDetails{
				ScheduledStartTime: "2024-01-01T00:00:00Z",
			}},
			want:    1704067200,
			wantErr: false,
		},
		"details is not nil, ScheduledStartTime is empty": {
			args: args{details: &youtube.VideoLiveStreamingDetails{
				ScheduledStartTime: "",
			}},
			want:    0,
			wantErr: true,
		},
		"details is not nil, ScheduledStartTime is invalid": {
			args: args{details: &youtube.VideoLiveStreamingDetails{
				ScheduledStartTime: "invalid",
			}},
			want:    0,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			// Arrange
			t.Parallel()
			// Act
			got, err := extractScheduledAtUnix(tt.args.details)
			// Assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
