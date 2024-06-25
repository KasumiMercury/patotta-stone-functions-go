package rss

import "time"

type Item struct {
	channelID       string
	sourceID        string
	title           string
	description     string
	publishedAtUnix int64
	updatedAtUnix   int64
}

func NewRssItem(channelID, sourceID, title, description string, publishedAtUnix, updatedAtUnix int64) *Item {
	return &Item{
		channelID:       channelID,
		sourceID:        sourceID,
		title:           title,
		description:     description,
		publishedAtUnix: publishedAtUnix,
		updatedAtUnix:   updatedAtUnix,
	}
}

func (r *Item) ChannelID() string {
	return r.channelID
}
func (r *Item) SourceID() string {
	return r.sourceID
}
func (r *Item) Title() string {
	return r.title
}
func (r *Item) Description() string {
	return r.description
}
func (r *Item) PublishedAtUnix() int64 {
	return r.publishedAtUnix
}
func (r *Item) PublishedAt() *time.Time {
	if r.publishedAtUnix == 0 {
		return nil
	}
	t := time.Unix(r.publishedAtUnix, 0)
	return &t
}
func (r *Item) UpdatedAtUnix() int64 {
	return r.updatedAtUnix
}
func (r *Item) UpdatedAt() *time.Time {
	if r.updatedAtUnix == 0 {
		return nil
	}
	t := time.Unix(r.updatedAtUnix, 0)
	return &t
}
