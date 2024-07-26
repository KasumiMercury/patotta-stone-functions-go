package service

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/config"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapters/api"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapters/db/realtime"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapters/rss"
	rssDto "github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapters/rss/dto"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/domain/video"
	"log/slog"
	"sort"
	"time"
)

var ytRssURL = "https://www.youtube.com/feeds/videos.xml?channel_id="

type SyncService struct {
	config  config.Config
	rssRepo rss.RSSRepository
	apiRepo api.ApiRepository
	rtdRepo realtime.RealtimeRepository
}

func NewSyncService(c config.Config, r rss.RSSRepository, a api.ApiRepository, rt realtime.RealtimeRepository) *SyncService {
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

	// Get updated videos from RSS
	rssItemList := make([]rssDto.Item, 0, 5)
	for _, c := range s.config.ChannelIDs() {
		// generate rss url
		url := ytRssURL + c
		// fetch rss items
		items, err := s.rssRepo.FetchRssItems(ctx, url, luu)
		if err != nil {
			return err
		}
		rssItemList = append(rssItemList, items...)
	}

	// Extract source IDs from updated rssItemList
	sidList := make([]string, 0, len(rssItemList))
	for _, r := range rssItemList {
		sidList = append(sidList, r.SourceID)
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

	// Update the video details in the database(RealtimeDB)
	videos := make([]video.Video, 0, len(vdList))

	for _, vd := range vdList {
		// merge video info and rss info
		m, err := video.NewVideo(
			vd.ChannelId,
			vd.Id,
			vd.Title,
			vd.Description,
			vd.ChatId,
			vd.Status,
			vd.PublishedAt.Unix(),
			vd.ScheduledAt.Unix(),
			time.Now().Unix(),
		)
		if err != nil {
			slog.Error(
				"Failed to create a video",
				"sourceID", vd.Id,
				"error", err,
			)
			continue
		}

		videos = append(videos, *m)
	}

	if len(videos) == 0 {
		slog.Info("No new videos found")
		return nil
	}

	// Sort the merged video info by updated time
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].UpdatedAtUnix() > videos[j].UpdatedAtUnix()
	})

	// Upsert the merged video info into the database(RealtimeDB)
	if err := s.rtdRepo.UpsertRecords(ctx, videos); err != nil {
		return err
	}

	return nil
}
