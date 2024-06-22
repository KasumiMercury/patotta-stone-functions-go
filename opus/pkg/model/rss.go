package model

import "time"

type Rss struct {
	ChannelID       string
	SourceID        string
	Title           string
	Description     string
	PublishedAtUnix int64
	UpdatedAtUnix   int64
}

func NewRss(channelID, sourceID, title, description string, publishedAtUnix, updatedAtUnix int64) *Rss {
	return &Rss{
		ChannelID:       channelID,
		SourceID:        sourceID,
		Title:           title,
		Description:     description,
		PublishedAtUnix: publishedAtUnix,
		UpdatedAtUnix:   updatedAtUnix,
	}
}

func (r *Rss) GetChannelID() string {
	return r.ChannelID
}
func (r *Rss) GetSourceID() string {
	return r.SourceID
}
func (r *Rss) GetTitle() string {
	return r.Title
}
func (r *Rss) GetDescription() string {
	return r.Description
}
func (r *Rss) GetPublishedAtUnix() int64 {
	return r.PublishedAtUnix
}
func (r *Rss) GetPublishedAt() time.Time {
	return time.Unix(r.PublishedAtUnix, 0)
}
func (r *Rss) GetUpdatedAtUnix() int64 {
	return r.UpdatedAtUnix
}
func (r *Rss) GetUpdatedAt() time.Time {
	return time.Unix(r.UpdatedAtUnix, 0)
}

type UpdatedItem struct {
	Record  VideoRecord
	RssItem Rss
}

func NewUpdatedItem(record VideoRecord, rssItem Rss) *UpdatedItem {
	return &UpdatedItem{
		Record:  record,
		RssItem: rssItem,
	}
}

func (ui *UpdatedItem) GetRecord() VideoRecord {
	return ui.Record
}
func (ui *UpdatedItem) GetRssItem() Rss {
	return ui.RssItem
}

type RssProcess struct {
	NewItems      []Rss
	TitleUpdated  []UpdatedItem
	DescUpdated   []UpdatedItem
	StatusUpdated []UpdatedItem
}

func NewRssProcess(newItems []Rss, titleUpdated, descUpdated, statusUpdated []UpdatedItem) *RssProcess {
	return &RssProcess{
		NewItems:      newItems,
		TitleUpdated:  titleUpdated,
		DescUpdated:   descUpdated,
		StatusUpdated: statusUpdated,
	}
}

func (rp *RssProcess) GetNewItems() []Rss {
	return rp.NewItems
}
func (rp *RssProcess) GetTitleUpdated() []UpdatedItem {
	return rp.TitleUpdated
}
func (rp *RssProcess) GetDescUpdated() []UpdatedItem {
	return rp.DescUpdated
}
func (rp *RssProcess) GetStatusUpdated() []UpdatedItem {
	return rp.StatusUpdated
}
