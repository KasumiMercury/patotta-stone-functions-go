package sync

import (
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/core/domain/api"
)

type VideoSyncItem struct {
	sourceID string
	content  *Content
	detail   *Detail
}

func NewItem(sourceID string, content *Content, detail *Detail) *VideoSyncItem {
	return &VideoSyncItem{
		sourceID: sourceID,
		content:  content,
		detail:   detail,
	}
}

type Content struct {
	sourceID      string
	title         string
	description   string
	updatedAtUnix int64
}

func NewContent(sourceID string, title string, description string, updatedAtUnix int64) *Content {
	return &Content{
		sourceID:      sourceID,
		title:         title,
		description:   description,
		updatedAtUnix: updatedAtUnix,
	}
}

type Detail struct {
	sourceID        string
	chatID          string
	status          api.Status
	publishedAtUnix int64
	scheduledAtUnix int64
}

func NewDetail(sourceID string, chatID string, status api.Status, publishedAtUnix int64, scheduledAtUnix int64) *Detail {
	return &Detail{
		sourceID:        sourceID,
		chatID:          chatID,
		status:          status,
		publishedAtUnix: publishedAtUnix,
		scheduledAtUnix: scheduledAtUnix,
	}
}
