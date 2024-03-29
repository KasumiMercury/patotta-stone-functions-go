package model

type Rss struct {
	ChannelID       string
	SourceID        string
	Title           string
	Description     string
	PublishedAtUnix int64
	UpdatedAtUnix   int64
}
