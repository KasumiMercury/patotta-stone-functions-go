package rss

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/rss"
	"github.com/mmcdole/gofeed"
)

type Client struct{}

func NewRssClient() *Client {
	return &Client{}
}

func (c *Client) GetRssItems(ctx context.Context, url string) ([]rss.Item, error) {
	feed, err := gofeed.NewParser().ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	items := make([]rss.Item, 0, len(feed.Items))
	for _, i := range feed.Items {
		item := rss.NewRssItem(
			i.Extensions["yt"]["channelId"][0].Value,
			i.Extensions["yt"]["videoId"][0].Value,
			i.Title,
			extractDescriptionFromRssItem(i),
			i.PublishedParsed.Unix(),
			i.UpdatedParsed.Unix(),
		)
		items = append(items, *item)
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
