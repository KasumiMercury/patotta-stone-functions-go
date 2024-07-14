package rss

import (
	"reflect"
	"testing"
	"time"
)

func TestItem_NillablePublishedAt(t *testing.T) {
	type fields struct {
		channelID       string
		sourceID        string
		title           string
		description     string
		publishedAtUnix int64
		updatedAtUnix   int64
	}
	tests := []struct {
		name   string
		fields fields
		want   *time.Time
	}{
		{
			name: "publishedAtUnix is 0",
			fields: fields{
				publishedAtUnix: 0,
			},
			want: nil,
		},
		{
			name: "publishedAtUnix is not 0",
			fields: fields{
				publishedAtUnix: 1610000000,
			},
			want: timePtr(time.Unix(1610000000, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Item{
				channelID:       tt.fields.channelID,
				sourceID:        tt.fields.sourceID,
				title:           tt.fields.title,
				description:     tt.fields.description,
				publishedAtUnix: tt.fields.publishedAtUnix,
				updatedAtUnix:   tt.fields.updatedAtUnix,
			}
			if got := r.NillablePublishedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillablePublishedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_NillableUpdatedAt(t *testing.T) {
	type fields struct {
		channelID       string
		sourceID        string
		title           string
		description     string
		publishedAtUnix int64
		updatedAtUnix   int64
	}
	tests := []struct {
		name   string
		fields fields
		want   *time.Time
	}{
		{
			name: "updatedAtUnix is 0",
			fields: fields{
				updatedAtUnix: 0,
			},
			want: nil,
		},
		{
			name: "updatedAtUnix is not 0",
			fields: fields{
				updatedAtUnix: 1610000000,
			},
			want: timePtr(time.Unix(1610000000, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Item{
				channelID:       tt.fields.channelID,
				sourceID:        tt.fields.sourceID,
				title:           tt.fields.title,
				description:     tt.fields.description,
				publishedAtUnix: tt.fields.publishedAtUnix,
				updatedAtUnix:   tt.fields.updatedAtUnix,
			}
			if got := r.NillableUpdatedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NillableUpdatedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
