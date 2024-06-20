package usecase

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/repository"
	"time"
)

type Video interface {
	GetVideoRecordMap(ctx context.Context, sourceIDs []string) (map[string]model.VideoRecord, error)
	SaveNewVideo(ctx context.Context, infos []model.VideoInfo, details map[string]model.VideoDetail) error
	UpdateScheduledAt(ctx context.Context, sourceID string, scheduledAt time.Time) error
	UpdateStatus(ctx context.Context, sourceID, status string) error
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

func (u *videoUsecase) SaveNewVideo(ctx context.Context, infos []model.VideoInfo, details map[string]model.VideoDetail) error {
	rec := make([]model.VideoRecord, 0, len(infos))
	for _, i := range infos {
		rec = append(rec, model.VideoRecord{
			SourceID:    i.SourceID,
			Title:       i.Title,
			Status:      details[i.SourceID].Status,
			ChatID:      details[i.SourceID].ChatID,
			ScheduledAt: time.Unix(details[i.SourceID].ScheduledAtUnix, 0),
			UpdatedAt:   time.Unix(i.UpdatedAtUnix, 0),
		})
	}

	if err := u.supaRepo.InsertVideoRecords(ctx, rec); err != nil {
		return err
	}

	return nil
}

func (u *videoUsecase) UpdateScheduledAt(ctx context.Context, sourceID string, scheduledAt time.Time) error {
	return u.supaRepo.UpdateScheduledAtBySourceID(ctx, sourceID, scheduledAt)
}

func (u *videoUsecase) UpdateStatus(ctx context.Context, sourceID, status string) error {
	return u.supaRepo.UpdateStatusBySourceID(ctx, sourceID, status)
}
