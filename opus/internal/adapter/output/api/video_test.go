package api

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/api"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/api/youtube/v3"
	"reflect"
	"testing"
)

func TestYouTubeVideo_FetchScheduledAtByVideoIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClient(ctrl)

	type args struct {
		videoIDs []string
	}

	t.Parallel()

	tests := map[string]struct {
		args      args
		mockSetup func(*mocks.MockClient)
		want      []api.LiveScheduleInfo
		wantErr   bool
	}{
		"success_single_video": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.LiveScheduleInfo{
				func() api.LiveScheduleInfo {
					l := api.NewLiveScheduleInfo("videoID")
					l.SetScheduledAtUnix(1704067200)
					return *l
				}(),
			},
			wantErr: false,
		},
		"success_multiple_videos": {
			args: args{videoIDs: []string{"videoID1", "videoID2"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID1", "videoID2"}), // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID1",
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
						{
							Id: "videoID2",
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.LiveScheduleInfo{
				func() api.LiveScheduleInfo {
					l := api.NewLiveScheduleInfo("videoID1")
					l.SetScheduledAtUnix(1704067200)
					return *l
				}(),
				func() api.LiveScheduleInfo {
					l := api.NewLiveScheduleInfo("videoID2")
					l.SetScheduledAtUnix(1704067200)
					return *l
				}(),
			},
			wantErr: false,
		},
		"error_api_call_fails": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),              // id
				).Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: true,
		},
		"error_scheduledStartTime_is_empty": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "",
							},
						},
					},
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		"error_scheduledStartTime_is_invalid": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "invalid",
							},
						},
					},
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			// Arrange
			tt.mockSetup(mockClient)

			c := &YouTubeVideo{
				clt: mockClient,
			}

			// Act
			got, err := c.FetchScheduledAtByVideoIDs(context.Background(), tt.args.videoIDs)
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

func TestYouTubeVideo_FetchVideoDetailsByVideoIDs(t *testing.T) {
	type fields struct {
		clt output.Client
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
