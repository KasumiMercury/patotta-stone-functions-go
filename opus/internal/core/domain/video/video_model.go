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

type Builder interface {
	SetChannelID(string) Builder
	SetSourceID(string) Builder
	SetTitle(string) Builder
	SetDescription(string) Builder
	SetChatID(string) Builder
	SetStatus(string) Builder
	SetPublishedAtUnix(int64) Builder
	SetScheduledAtUnix(int64) Builder
	SetUpdatedAtUnix(int64) Builder
	Build() *Video
}

type builder struct {
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

func NewVideoBuilder() Builder {
	return &builder{}
}

func (vb *builder) SetChannelID(channelID string) Builder {
	vb.channelID = channelID
	return vb
}

func (vb *builder) SetSourceID(sourceID string) Builder {
	vb.channelID = sourceID
	return vb
}

func (vb *builder) SetTitle(title string) Builder {
	vb.title = title
	return vb
}

func (vb *builder) SetDescription(description string) Builder {
	vb.description = description
	return vb
}

func (vb *builder) SetChatID(chatID string) Builder {
	vb.chatID = chatID
	return vb
}

func (vb *builder) SetStatus(status string) Builder {
	vb.status = status
	return vb
}

func (vb *builder) SetPublishedAtUnix(publishedAtUnix int64) Builder {
	vb.publishedAtUnix = publishedAtUnix
	return vb
}

func (vb *builder) SetScheduledAtUnix(scheduledAtUnix int64) Builder {
	vb.scheduledAtUnix = scheduledAtUnix
	return vb
}

func (vb *builder) SetUpdatedAtUnix(updatedAtUnix int64) Builder {
	vb.updatedAtUnix = updatedAtUnix
	return vb
}

func (vb *builder) Build() *Video {
	return &Video{
		channelID:       vb.channelID,
		sourceID:        vb.sourceID,
		title:           vb.title,
		description:     vb.description,
		chatID:          vb.chatID,
		status:          vb.status,
		publishedAtUnix: vb.publishedAtUnix,
		scheduledAtUnix: vb.scheduledAtUnix,
		updatedAtUnix:   vb.updatedAtUnix,
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
