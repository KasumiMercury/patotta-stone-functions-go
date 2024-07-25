package dto

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type Item struct {
	ChannelID   string
	SourceID    string
	Title       string
	Description string
	PublishedAt synchro.Time[tz.AsiaTokyo]
	UpdatedAt   synchro.Time[tz.AsiaTokyo]
}
