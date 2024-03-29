package usecase

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/repository"
	"log/slog"
)

type Rss interface {
	FetchUpdatedRssItemsEachOfChannels(ctx context.Context, target []string) ([]model.Rss, error)
}

type rssUsecase struct {
	rssRepo  repository.Rss
	supaRepo repository.Supabase
}

func NewRssUsecase(rssRepo repository.Rss, supaRepo repository.Supabase) Rss {
	return &rssUsecase{
		rssRepo:  rssRepo,
		supaRepo: supaRepo,
	}
}

func (u *rssUsecase) FetchUpdatedRssItemsEachOfChannels(ctx context.Context, target []string) ([]model.Rss, error) {
	threshold, err := u.supaRepo.GetLastUpdatedUnixOfVideo(ctx)
	if err != nil {
		slog.Error(
			"Failed to get the threshold of updated time of video",
			"error", err,
		)
		return nil, err
	}

	items := make([]model.Rss, 0)
	for _, t := range target {
		updated, err := u.rssRepo.FetchUpdatedRssItems(ctx, t, threshold)
		if err != nil {
			return nil, err
		}

		items = append(items, updated...)
	}

	return items, nil
}
