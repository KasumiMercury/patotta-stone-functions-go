package model

type Rss struct {
	ChannelID       string
	SourceID        string
	Title           string
	Description     string
	PublishedAtUnix int64
	UpdatedAtUnix   int64
}

type UpdatedItem struct {
	Record  VideoRecord
	RssItem Rss
}

type RssProcess struct {
	NewItems      []Rss
	TitleUpdated  []UpdatedItem
	DescUpdated   []UpdatedItem
	StatusUpdated []UpdatedItem
}
