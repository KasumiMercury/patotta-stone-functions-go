package api

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/api"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output/mocks"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/status"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/api/youtube/v3"
	"testing"
)

func TestYouTubeVideo_FetchVideoDetailsByVideoIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClient(ctrl)

	type args struct {
		videoIDs []string
	}

	tests := map[string]struct {
		args      args
		mockSetup func(*mocks.MockClient)
		want      []api.VideoDetail
		wantErr   bool
	}{
		"success_single_live_video": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "live",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Live)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_single_upcoming_video": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "upcoming",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Upcoming)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_single_archived_of_non_live_video": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "none",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Archived)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_single_archived_of_completed_live_video": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "completed",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Archived)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_multiple_videos_live_and_upcoming": {
			args: args{videoIDs: []string{"videoID1", "videoID2"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID1",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "live",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
						{
							Id: "videoID2",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "upcoming",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID1")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Live)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID2")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Upcoming)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_multiple_videos_live_and_archived_of_non_live": {
			args: args{videoIDs: []string{"videoID1", "videoID2"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID1",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "live",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
						{
							Id: "videoID2",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "none",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID1")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Live)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID2")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Archived)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_multiple_videos_live_and_archived_of_completed_live": {
			args: args{videoIDs: []string{"videoID1", "videoID2"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID1",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "live",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
						{
							Id: "videoID2",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "completed",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID1")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Live)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID2")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Archived)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_multiple_videos_upcoming_and_archived_of_non_live": {
			args: args{videoIDs: []string{"videoID1", "videoID2"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID1",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "upcoming",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
						{
							Id: "videoID2",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "none",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID1")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Upcoming)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID2")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Archived)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_multiple_videos_upcoming_and_archived_of_completed_live": {
			args: args{videoIDs: []string{"videoID1", "videoID2"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID1",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "upcoming",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
						{
							Id: "videoID2",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "completed",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID1")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Upcoming)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID2")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Archived)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_multiple_videos_archived_of_non_live_and_completed_live": {
			args: args{videoIDs: []string{"videoID1", "videoID2"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID1",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "none",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
						{
							Id: "videoID2",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "completed",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID1")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Archived)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID2")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Archived)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"success_multiple_videos_but_only_one_video_live": {
			args: args{videoIDs: []string{"videoID1", "videoID2"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID1",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "live",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: []api.VideoDetail{
				func() api.VideoDetail {
					vd := api.NewVideoDetail("videoID1")
					vd.SetPublishedAtUnix(1704067200)
					vd.SetStatus(status.Live)
					err := vd.SetScheduledAtUnix(1704067200)
					if err != nil {
						return api.VideoDetail{}
					}
					return *vd
				}(),
			},
			wantErr: false,
		},
		"error_api_call_fails": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: true,
		},
		"error_snippet_is_nil": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
						},
					},
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		"error_publishedAt_is_invalid": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt: "invalid",
							},
						},
					},
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		"error_liveBroadcastContent_is_invalid": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "invalid",
							},
						},
					},
				}, nil)
			},
			want:    nil,
			wantErr: true,
		},
		"error_scheduledStartTime_is_empty": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "live",
							},
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
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),                                           // id
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "live",
							},
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
			got, err := c.FetchVideoDetailsByVideoIDs(context.Background(), tt.args.videoIDs)
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
		"success_multiple_but_only_one_video": {
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
					},
				}, nil)
			},
			want: []api.LiveScheduleInfo{
				func() api.LiveScheduleInfo {
					l := api.NewLiveScheduleInfo("videoID1")
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
			t.Parallel()

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

func TestYouTubeVideo_FetchScheduledAtByVideoIDsAbnormally(t *testing.T) {
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
	}{
		"error_api_call_fails": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"liveStreamingDetails"}), // part
					gomock.Eq([]string{"videoID"}),              // id
				).Return(nil, assert.AnError)
			},
			want: nil,
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
			want: nil,
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
			want: nil,
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			tt.mockSetup(mockClient)

			c := &YouTubeVideo{
				clt: mockClient,
			}

			// Act
			got, err := c.FetchScheduledAtByVideoIDs(context.Background(), tt.args.videoIDs)
			// Assert
			assert.Error(t, err)
			assert.Equal(t, tt.want, got)
		})

	}
}

func TestExtractScheduledAtUnixSuccessfully(t *testing.T) {
	type args struct {
		details *youtube.VideoLiveStreamingDetails
	}

	t.Parallel()

	tests := map[string]struct {
		args args
		want int64
	}{
		"details is nil": {
			args: args{details: nil},
			want: 0,
		},
		"details is not nil": {
			args: args{details: &youtube.VideoLiveStreamingDetails{
				ScheduledStartTime: "2024-01-01T00:00:00Z",
			}},
			want: 1704067200,
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
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractScheduledAtUnixAbnormally(t *testing.T) {
	type args struct {
		details *youtube.VideoLiveStreamingDetails
	}

	t.Parallel()

	tests := map[string]struct {
		args args
		want int64
	}{
		"details is not nil, ScheduledStartTime is empty": {
			args: args{details: &youtube.VideoLiveStreamingDetails{
				ScheduledStartTime: "",
			}},
			want: 0,
		},
		"details is not nil, ScheduledStartTime is invalid": {
			args: args{details: &youtube.VideoLiveStreamingDetails{
				ScheduledStartTime: "invalid",
			}},
			want: 0,
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
			assert.Error(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
