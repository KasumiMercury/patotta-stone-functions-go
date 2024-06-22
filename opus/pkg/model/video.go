package model

import (
	"github.com/uptrace/bun"
	"time"
)

type VideoInfo struct {
	SourceID      string
	Title         string
	UpdatedAtUnix int64
}

func NewVideoInfo(sourceID string, title string, updatedAtUnix int64) *VideoInfo {
	return &VideoInfo{
		SourceID:      sourceID,
		Title:         title,
		UpdatedAtUnix: updatedAtUnix,
	}
}

func (vi *VideoInfo) GetSourceID() string {
	return vi.SourceID
}
func (vi *VideoInfo) GetTitle() string {
	return vi.Title
}
func (vi *VideoInfo) GetUpdatedAtUnix() int64 {
	return vi.UpdatedAtUnix
}
func (vi *VideoInfo) GetUpdatedAt() time.Time {
	return time.Unix(vi.UpdatedAtUnix, 0)
}

type VideoDetail struct {
	SourceID        string
	ChatID          string
	Status          string
	ScheduledAtUnix int64
	PublishedAtUnix int64
}

func NewVideoDetail(sourceID string, chatID string, status string, scheduledAtUnix int64, publishedAtUnix int64) *VideoDetail {
	return &VideoDetail{
		SourceID:        sourceID,
		ChatID:          chatID,
		Status:          status,
		ScheduledAtUnix: scheduledAtUnix,
		PublishedAtUnix: publishedAtUnix,
	}
}

func (vd *VideoDetail) GetSourceID() string {
	return vd.SourceID
}
func (vd *VideoDetail) GetChatID() string {
	return vd.ChatID
}
func (vd *VideoDetail) GetStatus() string {
	return vd.Status
}
func (vd *VideoDetail) GetScheduledAtUnix() int64 {
	return vd.ScheduledAtUnix
}
func (vd *VideoDetail) GetScheduledAt() time.Time {
	return time.Unix(vd.ScheduledAtUnix, 0)
}

type VideoRecord struct {
	bun.BaseModel `bun:"table:videos"`

	SourceID    string    `bun:",type:varchar(255)"`
	Title       string    `bun:",type:varchar(255)"`
	Status      string    `bun:",type:varchar(255)"`
	ChatID      string    `bun:",type:varchar(255)"`
	ScheduledAt time.Time `bun:",type:timestamptz"`
	UpdatedAt   time.Time `bun:",type:timestamptz"`
}

func NewVideoRecord(sourceID string, title string, status string, chatID string, scheduledAt time.Time, updatedAt time.Time) *VideoRecord {
	return &VideoRecord{
		SourceID:    sourceID,
		Title:       title,
		Status:      status,
		ChatID:      chatID,
		ScheduledAt: scheduledAt,
		UpdatedAt:   updatedAt,
	}
}

func (vr *VideoRecord) GetSourceID() string {
	return vr.SourceID
}
func (vr *VideoRecord) GetTitle() string {
	return vr.Title
}
func (vr *VideoRecord) GetStatus() string {
	return vr.Status
}
func (vr *VideoRecord) GetChatID() string {
	return vr.ChatID
}
func (vr *VideoRecord) GetScheduledAt() time.Time {
	return vr.ScheduledAt
}
func (vr *VideoRecord) GetUpdatedAt() time.Time {
	return vr.UpdatedAt
}

type VideoSchedule struct {
	SourceID    string
	Status      string
	ScheduledAt time.Time
}

func NewVideoSchedule(sourceID string, status string, scheduledAt time.Time) *VideoSchedule {
	return &VideoSchedule{
		SourceID:    sourceID,
		Status:      status,
		ScheduledAt: scheduledAt,
	}
}

func (vs *VideoSchedule) GetSourceID() string {
	return vs.SourceID
}
func (vs *VideoSchedule) GetStatus() string {
	return vs.Status
}
func (vs *VideoSchedule) GetScheduledAt() time.Time {
	return vs.ScheduledAt
}
