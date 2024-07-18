package api

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/status"
)

type DetailResponse struct {
	Id          string
	Title       string
	Description string
	Status      status.Status
	PublishedAt synchro.Time[tz.AsiaTokyo]
	ScheduledAt synchro.Time[tz.AsiaTokyo]
	ChatId      string
}
