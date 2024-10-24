package api

import (
	"context"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapters/api/dto"
	mocks "github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapters/youtube/mock"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/api/youtube/v3"
	"testing"
	"time"
)

func TestYouTubeVideo_FetchVideoDetailsByVideoIDsSuccessfully(t *testing.T) {
	t.Parallel()

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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID"}),                                           // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "live",
									ChannelId:            "channelID",
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
					Id:        "videoID",
					ChannelId: "channelID",
					Status:    status.Live,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID"}),                                           // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "upcoming",
									ChannelId:            "channelID",
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
					Id:        "videoID",
					ChannelId: "channelID",
					Status:    status.Upcoming,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID"}),                                           // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "none",
									ChannelId:            "channelID",
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
					Id:        "videoID",
					ChannelId: "channelID",
					Status:    status.Archived,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID"}),                                           // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "completed",
									ChannelId:            "channelID",
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
					Id:        "videoID",
					ChannelId: "channelID",
					Status:    status.Archived,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID1", "videoID2"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID1",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-02T00:00:00Z",
									LiveBroadcastContent: "live",
									ChannelId:            "channelID",
								},
								LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
									ScheduledStartTime: "2024-01-02T00:00:00Z",
								},
							},
							{
								Id: "videoID2",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "upcoming",
									ChannelId:            "channelID",
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
					Id:        "videoID1",
					ChannelId: "channelID",
					Status:    status.Live,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:        "videoID2",
					ChannelId: "channelID",
					Status:    status.Upcoming,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID1", "videoID2"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID1",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-02T00:00:00Z",
									LiveBroadcastContent: "live",
									ChannelId:            "channelID",
								},
								LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
									ScheduledStartTime: "2024-01-02T00:00:00Z",
								},
							},
							{
								Id: "videoID2",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "none",
									ChannelId:            "channelID",
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
					Id:        "videoID1",
					ChannelId: "channelID",
					Status:    status.Live,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:        "videoID2",
					ChannelId: "channelID",
					Status:    status.Archived,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID1", "videoID2"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID1",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-02T00:00:00Z",
									LiveBroadcastContent: "live",
									ChannelId:            "channelID",
								},
								LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
									ScheduledStartTime: "2024-01-02T00:00:00Z",
								},
							},
							{
								Id: "videoID2",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "completed",
									ChannelId:            "channelID",
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
					Id:        "videoID1",
					ChannelId: "channelID",
					Status:    status.Live,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:        "videoID2",
					ChannelId: "channelID",
					Status:    status.Archived,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID1", "videoID2"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID1",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-02T00:00:00Z",
									LiveBroadcastContent: "upcoming",
									ChannelId:            "channelID",
								},
								LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
									ScheduledStartTime: "2024-01-02T00:00:00Z",
								},
							},
							{
								Id: "videoID2",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "none",
									ChannelId:            "channelID",
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
					Id:        "videoID1",
					ChannelId: "channelID",
					Status:    status.Upcoming,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:        "videoID2",
					ChannelId: "channelID",
					Status:    status.Archived,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID1", "videoID2"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID1",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-02T00:00:00Z",
									LiveBroadcastContent: "upcoming",
									ChannelId:            "channelID",
								},
								LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
									ScheduledStartTime: "2024-01-02T00:00:00Z",
								},
							},
							{
								Id: "videoID2",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "completed",
									ChannelId:            "channelID",
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
					Id:        "videoID1",
					ChannelId: "channelID",
					Status:    status.Upcoming,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:        "videoID2",
					ChannelId: "channelID",
					Status:    status.Archived,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID1", "videoID2"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID1",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-02T00:00:00Z",
									LiveBroadcastContent: "none",
									ChannelId:            "channelID",
								},
								LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
									ScheduledStartTime: "2024-01-02T00:00:00Z",
								},
							},
							{
								Id: "videoID2",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "completed",
									ChannelId:            "channelID",
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
					Id:        "videoID1",
					ChannelId: "channelID",
					Status:    status.Archived,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				{
					Id:        "videoID2",
					ChannelId: "channelID",
					Status:    status.Archived,
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
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID1", "videoID2"}),                              // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"snippet", "contentDetails", "liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID1", "videoID2"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID1",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "live",
									ChannelId:            "channelID",
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
					Id:        "videoID1",
					ChannelId: "channelID",
					Status:    status.Live,
					PublishedAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"success_0_length_video_ids": {
			args: args{videoIDs: []string{}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			want: make([]dto.DetailResponse, 0),
		},
		"success_long_length_video_ids": {
			args: args{
				videoIDs: make([]string, 130),
			},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().
					VideoList(gomock.Any(), gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}), gomock.Any()).
					Times(3).
					DoAndReturn(func(_ context.Context, part []string, ids []string) (*youtube.VideoListResponse, error) {
						if len(ids) > 50 {
							return nil, assert.AnError
						}

						items := make([]*youtube.Video, 0, len(ids))
						for i := 0; i < len(ids); i++ {
							items = append(items, &youtube.Video{
								Id: "videoID",
								Snippet: &youtube.VideoSnippet{
									PublishedAt:          "2024-01-01T00:00:00Z",
									LiveBroadcastContent: "completed",
									ChannelId:            "channelID",
								},
								LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
									ScheduledStartTime: "2024-01-01T00:00:00Z",
								},
							})
						}

						return &youtube.VideoListResponse{
							Items: items,
						}, nil
					})
			},
			want: func() []dto.DetailResponse {
				var res []dto.DetailResponse
				for i := 0; i < 130; i++ {
					res = append(res, dto.DetailResponse{
						Id:        "videoID",
						ChannelId: "channelID",
						Status:    status.Archived,
						PublishedAt: synchro.In[tz.AsiaTokyo](
							time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
						ScheduledAt: synchro.In[tz.AsiaTokyo](
							time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					})
				}
				return res
			}(),
		},
		"abnormally_snippet_not_found": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}),
					gomock.Eq([]string{"videoID"}),
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
						},
					},
				}, nil)
			},
			want: make([]dto.DetailResponse, 0),
		},
		"abnormally_failed_to_parse_published_at": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}),
					gomock.Eq([]string{"videoID"}),
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
			want: make([]dto.DetailResponse, 0),
		},
		"abnormally_failed_to_match_video_status": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}),
					gomock.Eq([]string{"videoID"}),
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "unknown",
								ChannelId:            "channelID",
							},
						},
					},
				}, nil)
			},
			want: make([]dto.DetailResponse, 0),
		},
		"abnormally_failed_to_parse_scheduled_at": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}),
					gomock.Eq([]string{"videoID"}),
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt:          "2024-01-01T00:00:00Z",
								LiveBroadcastContent: "live",
								ChannelId:            "channelID",
							},
							LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
								ScheduledStartTime: "invalid",
							},
						},
					},
				}, nil)
			},
			want: make([]dto.DetailResponse, 0),
		},
		"abnormally_LiveBroadcastContent_not_found": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}),
					gomock.Eq([]string{"videoID"}),
				).Return(&youtube.VideoListResponse{
					Items: []*youtube.Video{
						{
							Id: "videoID",
							Snippet: &youtube.VideoSnippet{
								PublishedAt: "2024-01-01T00:00:00Z",
							},
						},
					},
				}, nil)
			},
			want: make([]dto.DetailResponse, 0),
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockClient(ctrl)
			tt.mockSetup(mockClient)

			c := &YouTubeVideo{
				clt: mockClient,
			}

			// Act
			got, err := c.FetchVideoDetailsByVideoIDs(context.Background(), tt.args.videoIDs)
			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !cmp.Equal(tt.want, got) {
				t.Errorf("unexpected response: %v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestYouTubeVideo_FetchVideoDetailsByVideoIDsError(t *testing.T) {
	t.Parallel()

	type args struct {
		videoIDs []string
	}

	tests := map[string]struct {
		args      args
		mockSetup func(*mocks.MockClient)
	}{
		"error_api_call_failed": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"snippet", "contentDetails", "liveStreamingDetails"}),
					gomock.Eq([]string{"videoID"}),
				).Return(nil, assert.AnError)
			},
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockClient(ctrl)
			tt.mockSetup(mockClient)

			c := &YouTubeVideo{
				clt: mockClient,
			}

			// Act
			_, err := c.FetchVideoDetailsByVideoIDs(context.Background(), tt.args.videoIDs)
			// Assert
			assert.Error(t, err)
		})
	}
}

func TestYouTubeVideo_FetchScheduledAtByVideoIDsSuccessfully(t *testing.T) {
	t.Parallel()

	type args struct {
		videoIDs []string
	}

	tests := map[string]struct {
		args      args
		mockSetup func(*mocks.MockClient)
		want      []dto.ScheduleResponse
	}{
		"success_single_live_video": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID"}),              // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID"}, ids)
					}).
					Return(&youtube.VideoListResponse{
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
			want: []dto.ScheduleResponse{
				{
					Id: "videoID",
					ScheduledAt: synchro.In[tz.AsiaTokyo](
						time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"success_non_live_video": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().
					VideoList(
						gomock.Any(),
						gomock.Eq([]string{"liveStreamingDetails"}), // part
						gomock.Eq([]string{"videoID"}),              // id
					).
					Times(1).
					Do(func(_ context.Context, part []string, ids []string) {
						assert.Equal(t, []string{"liveStreamingDetails"}, part)
						assert.Equal(t, []string{"videoID"}, ids)
					}).
					Return(&youtube.VideoListResponse{
						Items: []*youtube.Video{
							{
								Id: "videoID",
							},
						},
					}, nil)
			},
			want: []dto.ScheduleResponse{
				{
					Id:          "videoID",
					ScheduledAt: synchro.Time[tz.AsiaTokyo]{},
				},
			},
		},
		"success_0_length_video_ids": {
			args: args{videoIDs: []string{}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			want: make([]dto.ScheduleResponse, 0),
		},
		"success_long_length_video_ids": {
			args: args{
				videoIDs: make([]string, 130),
			},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().
					VideoList(gomock.Any(), gomock.Eq([]string{"liveStreamingDetails"}), gomock.Any()).
					Times(3).
					DoAndReturn(func(_ context.Context, part []string, ids []string) (*youtube.VideoListResponse, error) {
						if len(ids) > 50 {
							return nil, assert.AnError
						}

						items := make([]*youtube.Video, 0, len(ids))
						for i := 0; i < len(ids); i++ {
							items = append(items, &youtube.Video{
								Id: "videoID",
								LiveStreamingDetails: &youtube.VideoLiveStreamingDetails{
									ScheduledStartTime: "2024-01-01T00:00:00Z",
								},
							})
						}

						return &youtube.VideoListResponse{
							Items: items,
						}, nil
					})
			},
			want: func() []dto.ScheduleResponse {
				var res []dto.ScheduleResponse
				for i := 0; i < 130; i++ {
					res = append(res, dto.ScheduleResponse{
						Id: "videoID",
						ScheduledAt: synchro.In[tz.AsiaTokyo](
							time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					})
				}
				return res
			}(),
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockClient(ctrl)
			tt.mockSetup(mockClient)

			c := &YouTubeVideo{
				clt: mockClient,
			}

			// Act
			got, err := c.FetchScheduledAtByVideoIDs(context.Background(), tt.args.videoIDs)
			// Assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !cmp.Equal(tt.want, got) {
				t.Errorf("unexpected response: %v", cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestYouTubeVideo_FetchScheduledAtByVideoIDsError(t *testing.T) {
	t.Parallel()

	type args struct {
		videoIDs []string
	}

	tests := map[string]struct {
		args      args
		mockSetup func(*mocks.MockClient)
	}{
		"error_api_call_failed": {
			args: args{videoIDs: []string{"videoID"}},
			mockSetup: func(m *mocks.MockClient) {
				m.EXPECT().VideoList(
					gomock.Any(),
					gomock.Eq([]string{"liveStreamingDetails"}),
					gomock.Eq([]string{"videoID"}),
				).Return(nil, assert.AnError)
			},
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockClient(ctrl)
			tt.mockSetup(mockClient)

			c := &YouTubeVideo{
				clt: mockClient,
			}

			// Act
			_, err := c.FetchScheduledAtByVideoIDs(context.Background(), tt.args.videoIDs)
			// Assert
			assert.Error(t, err)
		})
	}
}
