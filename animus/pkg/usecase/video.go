package usecase

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/repository"
	"log/slog"
)

type Video interface {
	GetVideoInfosByStatusFromSupabase(ctx context.Context, status []string) (map[string][]model.VideoInfo, error)
}

type videoUsecase struct {
	supaRepo repository.Supabase
}

func NewVideoUsecase(supaRepo repository.Supabase) Video {
	return &videoUsecase{
		supaRepo: supaRepo,
	}
}

func (u *videoUsecase) GetVideoInfosByStatusFromSupabase(ctx context.Context, status []string) (map[string][]model.VideoInfo, error) {
	rec, err := u.supaRepo.GetVideoInfoByStatus(ctx, status)
	if err != nil {
		slog.Error("Failed to get video records",
			"status", status,
			slog.Group("error", err),
		)
	}

	if len(rec) == 0 {
		return nil, nil
	}

	videoInfos := make(map[string][]model.VideoInfo)
	for _, r := range rec {
		videoInfos[r.Status] = append(videoInfos[r.Status], model.VideoInfo{
			SourceID: r.SourceID,
			ChatID:   r.ChatID,
		})
	}

	return videoInfos, nil
}
