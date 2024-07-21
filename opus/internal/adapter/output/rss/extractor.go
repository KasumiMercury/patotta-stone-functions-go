package rss

import (
	"context"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapter/output/rss/dto"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/port/output"
	"github.com/mmcdole/gofeed"
)

type Client struct {
	parser output.ParserRepository
}

func NewRssClient(p output.ParserRepository) *Client {
	return &Client{parser: p}
}

func (c *Client) FetchRssItems(ctx context.Context, url string, limitUnix int64) ([]dto.Item, error) {
	feed, err := c.parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	items := make([]dto.Item, 0, len(feed.Items))
	for _, i := range feed.Items {
		// if updated is less than or equal to limitUnix, skip
		ut := i.UpdatedParsed.Unix()
		if ut <= limitUnix {
			continue
		}

		items = append(items, dto.Item{
			ChannelID:   i.Extensions["yt"]["channelId"][0].Value,
			SourceID:    i.Extensions["yt"]["videoId"][0].Value,
			Title:       i.Title,
			Description: extractDescriptionFromRssItem(i),
			PublishedAt: synchro.In[tz.AsiaTokyo](*i.PublishedParsed),
			UpdatedAt:   synchro.In[tz.AsiaTokyo](*i.UpdatedParsed),
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
