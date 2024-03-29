package usecase

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/repository"
	"time"
)

type Video interface {
	GetVideoRecordMap(ctx context.Context, sourceIDs []string) (map[string]model.VideoRecord, error)
	SaveNewVideo(ctx context.Context, videos []model.VideoInfo) error
}

type videoUsecase struct {
	supaRepo repository.Supabase
}

func NewVideoUsecase(supaRepo repository.Supabase) Video {
	return &videoUsecase{
		supaRepo: supaRepo,
	}
}

func (u *videoUsecase) GetVideoRecordMap(ctx context.Context, sourceIDs []string) (map[string]model.VideoRecord, error) {
	records, err := u.supaRepo.GetRecordsBySourceIDs(ctx, sourceIDs)
	if err != nil {
		return nil, err
	}

	recordMap := make(map[string]model.VideoRecord)
	for _, r := range records {
		recordMap[r.SourceID] = r
	}

	return recordMap, nil
}

func (u *videoUsecase) SaveNewVideo(ctx context.Context, infos []model.VideoInfo) error {
	// TODO: implement
	rec := make([]model.VideoRecord, 0, len(infos))
	for _, i := range infos {
		rec = append(rec, model.VideoRecord{
			SourceID:    i.SourceID,
			Title:       i.Title,
			Status:      i.Status,
			ChatID:      i.ChatID,
			ScheduledAt: time.Unix(i.ScheduledAtUnix, 0),
			UpdatedAt:   time.Unix(i.UpdatedAtUnix, 0),
		})
	}

	if err := u.supaRepo.InsertVideoRecords(ctx, rec); err != nil {
		return err
	}

	return nil
}
