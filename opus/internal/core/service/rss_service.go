package service

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output"
)

type RssService struct {
	rtdRepo output.RealtimeRepository
	rssRepo output.RSSRepository
}

func NewRssService(rtd *output.RealtimeRepository) *RssService {
	return &RssService{
		rtdRepo: *rtd,
	}
}

func (r *RssService) UpdateVideosFromRssItem(ctx context.Context) error {
	// Get the latest update time of the video info in the database(RealtimeDB)
	// To eliminate RSS that has already been confirmed
	lutu, err := r.rtdRepo.GetLastUpdatedUnixOfVideo(ctx)
	if err != nil {
		return err
	}

	duri := "https://www.youtube.com/feeds/videos.xml?channel_id=UCQ0UDLQCjY0rmuxCDE38FGg"

	// Get updated videos from RSS
	rss, err := r.rssRepo.FetchRssItems(ctx, duri, lutu)
	if err != nil {
		return err
	}

	// Extract source IDs from updated rss
	sourceIDs := make([]string, 0, len(rss))
	for _, r := range rss {
		sourceIDs = append(sourceIDs, r.SourceID())
	}

	// TODO: Implement
	// Get video info of updated videos from YouTube Data API
	// TODO: Implement
	// Update video info
	// TODO: Implement

	return nil
}
