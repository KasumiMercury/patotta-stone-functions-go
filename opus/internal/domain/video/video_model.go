package video

import (
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
)

type Video struct {
	channelID       string
	sourceID        string
	title           string
	description     string
	chatID          string
	status          status.Status
	publishedAtUnix int64
	scheduledAtUnix int64
	updatedAtUnix   int64
}

func NewVideo(channelID, sourceID, title, description, chatID string, status status.Status, publishedAtUnix, scheduledAtUnix, updatedAtUnix int64) *Video {
	return &Video{
		channelID:       channelID,
		sourceID:        sourceID,
		title:           title,
		description:     description,
		chatID:          chatID,
		status:          status,
		publishedAtUnix: publishedAtUnix,
		scheduledAtUnix: scheduledAtUnix,
		updatedAtUnix:   updatedAtUnix,
	}
}
