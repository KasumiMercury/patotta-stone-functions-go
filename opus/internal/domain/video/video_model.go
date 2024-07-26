package video

import (
	"fmt"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
	"time"
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
	if v.scheduledAtUnix != 0 && v.scheduledAtUnix < v.publishedAtUnix {
		return fmt.Errorf("scheduledAtUnix must be greater than publishedAtUnix")
	}

	if v.updatedAtUnix == 0 {
		return fmt.Errorf("updatedAtUnix is required")
	}

	return nil
}

func (v *Video) ChannelID() string {
	return v.channelID
}
func (v *Video) SourceID() string {
	return v.sourceID
}
func (v *Video) Title() string {
	return v.title
}
func (v *Video) Description() string {
	return v.description
}
func (v *Video) ChatID() string {
	return v.chatID
}
func (v *Video) Status() status.Status {
	return v.status
}
func (v *Video) PublishedAtUnix() int64 {
	return v.publishedAtUnix
}
func (v *Video) PublishedAtUTC() time.Time {
	return time.Unix(v.publishedAtUnix, 0).UTC()
}
func (v *Video) ScheduledAtUnix() int64 {
	return v.scheduledAtUnix
}
func (v *Video) NillableScheduledAtUTC() *time.Time {
	if v.scheduledAtUnix == 0 {
		return nil
	}

	t := time.Unix(v.scheduledAtUnix, 0).UTC()
	return &t
}
func (v *Video) UpdatedAtUnix() int64 {
	return v.updatedAtUnix
}
func (v *Video) UpdatedAtUTC() time.Time {
	return time.Unix(v.updatedAtUnix, 0).UTC()
}
