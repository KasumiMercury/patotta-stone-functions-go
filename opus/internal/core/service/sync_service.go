package service

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/rss"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/video"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/config"
	"log/slog"
)

type SyncService struct {
	config  config.Config
	rssRepo output.RSSRepository
	apiRepo output.ApiRepository
	rtdRepo output.RealtimeRepository
}

func NewSyncService(c config.Config, r output.RSSRepository, a output.ApiRepository, rt output.RealtimeRepository) *SyncService {
	return &SyncService{
		config:  c,
		rssRepo: r,
		apiRepo: a,
		rtdRepo: rt,
	}
}

func (s *SyncService) SyncVideosWithRSS(ctx context.Context) error {
	// Get the latest update time of the video info in the database(RealtimeDB)
	// To eliminate RSS that has already been confirmed
	luu, err := s.rtdRepo.GetLastUpdatedUnixOfVideo(ctx)
	if err != nil {
		return err
	}

	duri := "https://www.youtube.com/feeds/videos.xml?channel_id=UCeLzT-7b2PBcunJplmWtoDg"

	// Get updated videos from RSS
	rssItemList, err := s.rssRepo.FetchRssItems(ctx, duri, luu)
	if err != nil {
		return err
	}

	// Extract source IDs from updated rssItemList
	sidList := make([]string, 0, len(rssItemList))
	for _, r := range rssItemList {
		sidList = append(sidList, r.SourceID())
	}

	// Get video details of updated videos from YouTube Data API
	vdList, err := s.apiRepo.FetchVideoDetailsByVideoIDs(ctx, sidList)
	if err != nil {
		return err
	}

	// if the difference between len(rssItemList) and len(vdList) is not 0, log it as a warning
	if len(rssItemList) != len(vdList) {
		slog.Warn(
			"Failed to get video details for all updated videos",
			"rssItemList", len(rssItemList),
			"videoDetails", len(vdList),
		)
	}

	// make an item map from rssItemList
	rssMap := make(map[string]rss.Item, len(rssItemList))
	for _, r := range rssItemList {
		rssMap[r.SourceID()] = r
	}

	// Update the video details in the database(RealtimeDB)
	videos := make([]video.Video, 0, len(vdList))

	for _, vd := range vdList {
		r, ok := rssMap[vd.SourceID()]
		if !ok {
			slog.Warn(
				"Video  not found in RSS",
				"sourceID", vd.SourceID(),
			)
			continue
		}

		// merge video info and rss info
		m := video.NewVideoBuilder(vd.SourceID()).
			SetChannelID(r.ChannelID()).
			SetTitle(r.Title()).
			SetDescription(r.Description()).
			SetChatID(vd.ChatID()).
			SetPublishedAtUnix(vd.PublishedAtUnix()).
			SetScheduledAtUnix(vd.ScheduledAtUnix()).
			SetUpdatedAtUnix(r.UpdatedAtUnix()).
			Build()

		videos = append(videos, *m)
	}

	if len(videos) == 0 {
		slog.Info("No new videos found")
		return nil
	}

	// Upsert the merged video info into the database(RealtimeDB)
	if err := s.rtdRepo.UpsertRecords(ctx, videos); err != nil {
		return err
	}

	return nil
}
