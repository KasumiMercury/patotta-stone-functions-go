package video

import (
	"fmt"
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

func NewVideo(channelID, sourceID, title, description, chatID string, status status.Status, publishedAtUnix, scheduledAtUnix, updatedAtUnix int64) (*Video, error) {
	v := &Video{
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

	if err := v.validate(); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Video) validate() error {
	if v.channelID == "" {
		return fmt.Errorf("channelID is required")
	}
	if v.sourceID == "" {
		return fmt.Errorf("sourceID is required")
	}
	if v.title == "" {
		return fmt.Errorf("title is required")
	}

	// description is optional
	// chatID is optional

	if v.status == status.Undefined {
		return fmt.Errorf("status is undefined")
	}

	if v.publishedAtUnix == 0 {
		return fmt.Errorf("publishedAtUnix is required")
	}

	// scheduledAtUnix is optional
	// scheduledAtUnix is must be greater than publishedAtUnix
	if v.scheduledAtUnix != 0 && v.scheduledAtUnix <= v.publishedAtUnix {
		return fmt.Errorf("scheduledAtUnix must be greater than publishedAtUnix")
	}

	if v.updatedAtUnix == 0 {
		return fmt.Errorf("updatedAtUnix is required")
	}

	return nil
}

// TODO: getter or convert to DTO
