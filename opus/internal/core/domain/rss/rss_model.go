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

func (r *Item) GetChannelID() string {
	return r.channelID
}
func (r *Item) GetSourceID() string {
	return r.sourceID
}
func (r *Item) GetTitle() string {
	return r.title
}
func (r *Item) GetDescription() string {
	return r.description
}
func (r *Item) GetPublishedAtUnix() int64 {
	return r.publishedAtUnix
}
func (r *Item) GetPublishedAt() *time.Time {
	if r.publishedAtUnix == 0 {
		return nil
	}
	t := time.Unix(r.publishedAtUnix, 0)
	return &t
}
func (r *Item) GetUpdatedAtUnix() int64 {
	return r.updatedAtUnix
}
func (r *Item) GetUpdatedAt() *time.Time {
	if r.updatedAtUnix == 0 {
		return nil
	}
	t := time.Unix(r.updatedAtUnix, 0)
	return &t
}
