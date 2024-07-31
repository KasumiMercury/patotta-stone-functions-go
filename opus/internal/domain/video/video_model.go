package video

import (
	"fmt"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
)

type Video struct {
	channelID   string
	sourceID    string
	title       string
	description string
	chatID      string
	status      status.Status
	publishedAt synchro.Time[tz.AsiaTokyo]
	scheduledAt synchro.Time[tz.AsiaTokyo]
	updatedAt   synchro.Time[tz.AsiaTokyo]
}

func NewVideo(channelID, sourceID, title, description, chatID string, status status.Status, publishedAt, scheduledAt, updatedAt synchro.Time[tz.AsiaTokyo]) (*Video, error) {
	v := &Video{
		channelID:   channelID,
		sourceID:    sourceID,
		title:       title,
		description: description,
		chatID:      chatID,
		status:      status,
		publishedAt: publishedAt,
		scheduledAt: scheduledAt,
		updatedAt:   updatedAt,
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

	if v.publishedAt.IsZero() {
		return fmt.Errorf("publishedAt is required")
	}

	// scheduledAtUnix is optional
	// scheduledAtUnix is must be greater than publishedAtUnix
	if !v.scheduledAt.IsZero() && v.scheduledAt.Before(v.publishedAt) {
		return fmt.Errorf("scheduledAtUnix must be greater than publishedAtUnix")
	}

	if v.updatedAt.IsZero() {
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
func (v *Video) PublishedAt() synchro.Time[tz.AsiaTokyo] {
	return v.publishedAt
}
func (v *Video) ScheduledAt() synchro.Time[tz.AsiaTokyo] {
	return v.scheduledAt
}
func (v *Video) UpdatedAt() synchro.Time[tz.AsiaTokyo] {
	return v.updatedAt
}
