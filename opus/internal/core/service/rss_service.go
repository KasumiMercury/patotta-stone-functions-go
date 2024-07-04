package service

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapter/output/db/realtime"
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

	// Get video records from RealtimeDB
	// for judging whether to newly register or update the video
	vrList, err := r.rtdRepo.GetRecordsBySourceIDs(ctx, sidList)
	if err != nil {
		return err
	}

	// make vrList map
	vrMap := make(map[string]*realtime.Record)
	for _, v := range vrList {
		vrMap[v.SourceID] = v
	}

	// Compare RSS items and video records
	// and divide them into new items and updated items

	newItems := make([]*video.Video, 0, len(rssItemList)/2)
	updatedItems := make([]*video.Video, 0, len(rssItemList))

	for _, v := range vdList {
		ri, ok := rssMap[v.SourceID()]
		if !ok {
			continue
		}

		// merge video info and rss info
		vr := video.NewVideo(
			ri.ChannelID(),
			v.SourceID(),
			ri.Title(),
			ri.Description(),
			v.ChatID(),
			v.Status(),
			v.PublishedAtUnix(),
			v.ScheduledAtUnix(),
			ri.UpdatedAtUnix(),
		)

		// If the source ID is not in the video records, it is a new item
		if _, ok := vrMap[v.SourceID()]; !ok {
			newItems = append(newItems, vr)
			continue
		} else {
			// If the source ID is in the video records, it is an updated item
			updatedItems = append(updatedItems, vr)
		}
	}

	// Update video info
	// TODO: Implement

	return nil
}
