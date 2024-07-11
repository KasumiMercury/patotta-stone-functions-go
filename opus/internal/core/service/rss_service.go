package service

import (
	"context"
	rssDomain "github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/rss"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/video"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output"
	"log/slog"
)

type RssService struct {
	rtdRepo output.RealtimeRepository
	rssRepo output.RSSRepository
	apiRepo output.ApiRepository
}

func NewRssService(rtd *output.RealtimeRepository, rss *output.RSSRepository, api *output.ApiRepository) *RssService {
	return &RssService{
		rtdRepo: *rtd,
		rssRepo: *rss,
		apiRepo: *api,
	}
}

func (r *RssService) UpdateVideosFromRssItem(ctx context.Context) error {
	// Get the latest update time of the video info in the database(RealtimeDB)
	// To eliminate RSS that has already been confirmed
	luu, err := r.rtdRepo.GetLastUpdatedUnixOfVideo(ctx)
	if err != nil {
		return err
	}

	duri := "https://www.youtube.com/feeds/videos.xml?channel_id=UCeLzT-7b2PBcunJplmWtoDg"

	// Get updated videos from RSS
	rssItemList, err := r.rssRepo.FetchRssItems(ctx, duri, luu)
	if err != nil {
		return err
	}

	// Extract source IDs from updated rssItemList
	sidList := make([]string, 0, len(rssItemList))
	for _, r := range rssItemList {
		sidList = append(sidList, r.SourceID())
	}

	// Get video details of updated videos from YouTube Data API
	vdList, err := r.apiRepo.FetchVideoDetailsByVideoIDs(ctx, sidList)
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

	// make rssItemList map
	rssMap := make(map[string]rssDomain.Item)
	for _, r := range rssItemList {
		rssMap[r.SourceID()] = r
	}

	videos := make([]video.Video, 0, len(vdList))
	for _, v := range vdList {
		ri, ok := rssMap[v.SourceID()]
		if !ok {
			continue
		}

		// merge video info and rss info
		m := video.NewVideoBuilder().
			SetChannelID(ri.ChannelID()).
			SetSourceID(ri.SourceID()).
			SetTitle(ri.Title()).
			SetStatus(v.Status()).
			SetChatID(v.ChatID()).
			SetPublishedAtUnix(v.PublishedAtUnix()).
			SetScheduledAtUnix(v.ScheduledAtUnix()).
			SetUpdatedAtUnix(ri.UpdatedAtUnix()).
			Build()

		videos = append(videos, *m)
	}

	if len(videos) == 0 {
		slog.Info("No new videos found")
		return nil
	}

	// Upsert the merged video info into the database(RealtimeDB)
	if err := r.rtdRepo.UpsertRecords(ctx, videos); err != nil {
		return err
	}

	return nil
}
