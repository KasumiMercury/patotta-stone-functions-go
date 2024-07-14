package video

import (
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/status"
	"time"
)

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

type Builder interface {
	SetChannelID(string) Builder
	SetTitle(string) Builder
	SetDescription(string) Builder
	SetChatID(string) Builder
	SetStatus(string) Builder
	SetPublishedAtUnix(int64) Builder
	SetScheduledAtUnix(int64) Builder
	SetUpdatedAtUnix(int64) Builder
	Build() *Video
}

type builderImpl struct {
	video *Video
}

func NewVideoBuilder(sourceID string) Builder {
	return &builderImpl{
		video: &Video{
			sourceID: sourceID,
			status:   status.Undefined.String(),
		},
	}
}

func (vb *builderImpl) SetChannelID(channelID string) Builder {
	vb.video.channelID = channelID
	return vb
}

func (vb *builderImpl) SetTitle(title string) Builder {
	vb.video.title = title
	return vb
}

func (vb *builderImpl) SetDescription(description string) Builder {
	vb.video.description = description
	return vb
}

func (vb *builderImpl) SetChatID(chatID string) Builder {
	vb.video.chatID = chatID
	return vb
}

func (vb *builderImpl) SetStatus(status string) Builder {
	vb.video.status = status
	return vb
}

func (vb *builderImpl) SetPublishedAtUnix(publishedAtUnix int64) Builder {
	vb.video.publishedAtUnix = publishedAtUnix
	return vb
}

func (vb *builderImpl) SetScheduledAtUnix(scheduledAtUnix int64) Builder {
	vb.video.scheduledAtUnix = scheduledAtUnix
	return vb
}

func (vb *builderImpl) SetUpdatedAtUnix(updatedAtUnix int64) Builder {
	vb.video.updatedAtUnix = updatedAtUnix
	return vb
}

func (vb *builderImpl) Build() *Video {
	return vb.video
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
