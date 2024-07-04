package service

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output"
	"log/slog"
)

type RssService struct {
	rtdRepo output.RealtimeRepository
	rssRepo output.RSSRepository
	apiRepo output.ApiRepository
}

func NewRssService(rtd *output.RealtimeRepository) *RssService {
	return &RssService{
		rtdRepo: *rtd,
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
	rss, err := r.rssRepo.FetchRssItems(ctx, duri, luu)
	if err != nil {
		return err
	}

	// Extract source IDs from updated rss
	sidList := make([]string, 0, len(rss))
	for _, r := range rss {
		sidList = append(sidList, r.SourceID())
	}

	// Get video details of updated videos from YouTube Data API
	vdList, err := r.apiRepo.FetchVideoDetailsByVideoIDs(ctx, sidList)
	if err != nil {
		return err
	}

	// if difference between len(rss) and len(vdList) is not 0, log it as a warning
	if len(rss) != len(vdList) {
		slog.Warn(
			"Failed to get video details for all updated videos",
			"rss", len(rss),
			"videoDetails", len(vdList),
		)
	}

	// Update video info
	// TODO: Implement

	return nil
}
