package usecase

import (
	"context"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/repository"
)

type Video interface {
	GetVideoRecordMap(ctx context.Context, sourceIDs []string) (map[string]model.VideoRecord, error)
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
