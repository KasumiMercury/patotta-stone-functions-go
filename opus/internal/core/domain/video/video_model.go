package video

import "time"

type Video struct {
	channelID       string
	sourceID        string
	title           string
	description     string
	chatID          string
	status          string
	publishedAtUnix int64
	scheduledAtUnix int64
	updatedAtUnix   int64
}

func NewVideo(channelID, sourceID, title, description string, chatID string, status string, publishedAtUnix, scheduledAtUnix, updatedAtUnix int64) *Video {
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
func (v *Video) NillableChatID() *string {
	if v.chatID == "" {
		return nil
	}
	return &v.chatID
}
func (v *Video) Status() string {
	return v.status
}
func (v *Video) PublishedAtUnix() int64 {
	return v.publishedAtUnix
}
func (v *Video) NillablePublishedAt() *time.Time {
	if v.publishedAtUnix == 0 {
		return nil
	}
	t := time.Unix(v.publishedAtUnix, 0)
	return &t
}
func (v *Video) ScheduledAtUnix() int64 {
	return v.scheduledAtUnix
}
func (v *Video) NillableScheduledAt() *time.Time {
	if v.scheduledAtUnix == 0 {
		return nil
	}
	t := time.Unix(v.scheduledAtUnix, 0)
	return &t
}
func (v *Video) UpdatedAtUnix() int64 {
	return v.updatedAtUnix
}
func (v *Video) NillableUpdatedAt() *time.Time {
	if v.updatedAtUnix == 0 {
		return nil
	}
	t := time.Unix(v.updatedAtUnix, 0)
	return &t
}
