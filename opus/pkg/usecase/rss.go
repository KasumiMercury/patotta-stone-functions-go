package usecase

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/repository"
	"log/slog"
)

type Rss interface {
	FetchUpdatedRssItemsEachOfChannels(ctx context.Context, target []string) error
}

type rssUsecase struct {
	rssRepo repository.Rss
}

func NewRssUsecase(rssRepo repository.Rss) Rss {
	return &rssUsecase{
		rssRepo: rssRepo,
	}
}

func (u *rssUsecase) FetchUpdatedRssItemsEachOfChannels(ctx context.Context, target []string) error {
	for _, t := range target {
		items, err := u.rssRepo.FetchUpdatedRssItems(ctx, t, 0)
		if err != nil {
			return err
		}

		slog.Debug("Fetched RSS items", slog.Group("channel", t, "items", items))
	}

	return nil
}
