package api

import (
	"context"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapter/output/api/dto"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output/mocks"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/status"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/api/youtube/v3"
	"testing"
	"time"
)

func TestYouTubeVideo_FetchVideoDetailsByVideoIDsSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClient(ctrl)

	type args struct {
		videoIDs []string
	}

	tests := map[string]struct {
		args      args
		mockSetup func(*mocks.MockClient)
		want      []dto.DetailResponse
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID",
					Status: status.Live,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID",
					Status: status.Upcoming,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID",
					Status: status.Archived,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			//want: []api.VideoDetail{
			//	func() api.VideoDetail {
			//		vd, _ := api.NewVideoDetail("videoID", "", status.Archived, 1704067200, 1704067200)
			//		return *vd
			//	}(),
			//},
			want: []dto.DetailResponse{
				{
					Id:     "videoID",
					Status: status.Archived,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID1",
					Status: status.Live,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:     "videoID2",
					Status: status.Upcoming,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID1",
					Status: status.Live,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:     "videoID2",
					Status: status.Archived,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID1",
					Status: status.Live,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:     "videoID2",
					Status: status.Archived,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID1",
					Status: status.Upcoming,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:     "videoID2",
					Status: status.Archived,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID1",
					Status: status.Upcoming,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:     "videoID2",
					Status: status.Archived,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID1",
					Status: status.Archived,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:     "videoID2",
					Status: status.Archived,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			want: []dto.DetailResponse{
				{
					Id:     "videoID1",
					Status: status.Live,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
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
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
