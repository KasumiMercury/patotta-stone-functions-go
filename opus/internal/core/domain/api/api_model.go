package api

import (
	"fmt"
	"log/slog"
	"time"
)

type VideoDetail struct {
	sourceID        string
	chatID          string
	status          Status
	publishedAtUnix int64
	scheduledAtUnix int64
}

func NewVideoDetail(sourceID string) *VideoDetail {
	return &VideoDetail{
		sourceID: sourceID,
		status:   Undefined,
	}
}

func (vd *VideoDetail) SetChatID(chatID string) {
	vd.chatID = chatID
}
func (vd *VideoDetail) SetStatus(status Status) {
	vd.status = status
}
func (vd *VideoDetail) SetPublishedAtUnix(publishedAtUnix int64) {
	vd.publishedAtUnix = publishedAtUnix
}
func (vd *VideoDetail) SetScheduledAtUnix(scheduledAtUnix int64) error {
	if vd.publishedAtUnix == 0 {
		slog.Warn(
			"publishedAtUnix is not set",
			"sourceID", vd.sourceID,
		)
		return fmt.Errorf("publishedAtUnix is not set for sourceID: %s", vd.sourceID)
	}
	vd.scheduledAtUnix = scheduledAtUnix
	return nil
}

func (vd *VideoDetail) GetSourceID() string {
	return vd.sourceID
}
func (vd *VideoDetail) GetNillableChatID() *string {
	if vd.chatID == "" {
		return nil
	}
	return &vd.chatID
}
func (vd *VideoDetail) GetStatus() string {
	return vd.status.String()
}
func (vd *VideoDetail) GetScheduledAtUnix() int64 {
	return vd.scheduledAtUnix
}
func (vd *VideoDetail) GetNillableScheduledAt() *time.Time {
	if vd.scheduledAtUnix == 0 {
		return nil
	}
	t := time.Unix(vd.scheduledAtUnix, 0)
	return &t
}
func (vd *VideoDetail) GetPublishedAtUnix() int64 {
	return vd.publishedAtUnix
}
func (vd *VideoDetail) GetNillablePublishedAt() *time.Time {
	if vd.publishedAtUnix == 0 {
		return nil
	}
	t := time.Unix(vd.publishedAtUnix, 0)
	return &t
}

type LiveScheduleInfo struct {
	sourceID        string
	scheduledAtUnix int64
}

func NewLiveScheduleInfo(sourceID string) *LiveScheduleInfo {
	return &LiveScheduleInfo{
		sourceID: sourceID,
	}
}

func (lsi *LiveScheduleInfo) SetScheduledAtUnix(scheduledAtUnix int64) {
	lsi.scheduledAtUnix = scheduledAtUnix
}

func (lsi *LiveScheduleInfo) GetSourceID() string {
	return lsi.sourceID
}
func (lsi *LiveScheduleInfo) GetScheduledAtUnix() int64 {
	return lsi.scheduledAtUnix
}
func (lsi *LiveScheduleInfo) GetNillableScheduledAt() *time.Time {
	if lsi.scheduledAtUnix == 0 {
		return nil
	}
	t := time.Unix(lsi.scheduledAtUnix, 0)
	return &t
}