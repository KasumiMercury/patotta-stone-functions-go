package infra

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
	"github.com/mmcdole/gofeed"
	"log/slog"
)

type RssRepository struct{}

func NewRssRepository() *RssRepository {
	return &RssRepository{}
}

func (r *RssRepository) FetchUpdatedRssItems(ctx context.Context, url string, threshold int64) ([]model.Rss, error) {
	feed, err := gofeed.NewParser().ParseURLWithContext(url, ctx)
	if err != nil {
		slog.Error(
			"failed to parse rss",
			slog.Group("fetchRss", "url", url,
				slog.Group("gofeed", "error", err),
			),
		)
		return nil, err
	}

	items := make([]model.Rss, 0, len(feed.Items))
	for _, i := range feed.Items {
		updated := i.UpdatedParsed.Unix()
		if updated <= threshold {
			continue
		}
		items = append(items, model.Rss{
			ChannelID:       i.Extensions["yt"]["channelId"][0].Value,
			SourceID:        i.Extensions["yt"]["videoId"][0].Value,
			Title:           i.Title,
			Description:     extractDescriptionFromRssItem(i),
			PublishedAtUnix: i.PublishedParsed.Unix(),
			UpdatedAtUnix:   updated,
		})
	}

	return items, nil
}

func extractDescriptionFromRssItem(i *gofeed.Item) string {
	if i == nil {
		return ""
	}

	if media, ok := i.Extensions["media"]; ok {
		if groups, ok := media["group"]; ok && len(groups) > 0 {
			group := groups[0].Children
			if description, ok := group["description"]; ok && len(description) > 0 {
				return description[0].Value
			}
		}
	}

	return ""
}
