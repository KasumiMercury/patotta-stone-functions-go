package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/repository"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/service"
	"log/slog"
	"os"
	"slices"
)

type Chat interface {
	FetchChatsFromStaticTargetVideo(ctx context.Context) ([]model.YTChat, error)
	FetchChatsFromUpcomingTargetVideo(ctx context.Context, upc []model.VideoInfo) ([]model.YTChat, error)
}

type chatUsecase struct {
	targetChannel []string
	chatSvc       service.Chat
	supaRepo      repository.Supabase
}

func NewChatUsecase(targetChannel []string, chatSvc service.Chat, supaRepo repository.Supabase) Chat {
	return &chatUsecase{
		targetChannel: targetChannel,
		chatSvc:       chatSvc,
		supaRepo:      supaRepo,
	}
}

func (u *chatUsecase) FetchChatsFromStaticTargetVideo(ctx context.Context) ([]model.YTChat, error) {
	// load info of static target video from environment variables
	stcEnv := os.Getenv("STATIC_TARGET")
	var stc model.VideoInfo
	if err := json.Unmarshal([]byte(stcEnv), &stc); err != nil {
		slog.Error("Failed to unmarshal STATIC_TARGET", "error", err)
		// If the environment variable is not set correctly, the function will panic.
		// (To prevent retries by CloudScheduler, the function should panic without returning error responses.)
		panic(fmt.Sprintf("Failed to unmarshal static target: %v", err))
	}

	// Fetch chats from the static target video
	stcChats, err := u.chatSvc.FetchChatsByVideoInfo(ctx, stc, 0)
	if err != nil {
		slog.Error("Failed to fetch chats from the static target video", slog.Group("staticTarget", "error", err))
		return nil, err
	}

	// Filter chats by author channel
	targetChats, _ := filterChatsByAuthorChannel(stcChats, u.targetChannel)

	if len(targetChats) == 0 {
		slog.Info("No chats from the static target video")
		return nil, nil
	}

	// Filter chats by the publishedAt
	newChats, err := u.filterChatsByPublishedAt(ctx, targetChats, stc.SourceID)
	if err != nil {
		slog.Error("Failed to filter chats by the publishedAt",
			slog.Group("staticTarget", "error", err),
		)
		return nil, err
	}

	// debug log
	slog.Debug("Fetched chats from the static target video", "count", len(newChats))

	return newChats, nil
}

func (u *chatUsecase) FetchChatsFromUpcomingTargetVideo(ctx context.Context, upc []model.VideoInfo) ([]model.YTChat, error) {
	// Fetch chats from the upcoming target video
	var video model.VideoInfo

	// If there are multiple videos with Upcoming status, calculate priority
	// the priority is based on fetched history
	switch len(upc) {
	case 0:
		slog.Info("No upcoming target video")
		return nil, nil
	case 1:
		video = upc[0]
	default:
		// Get fetched history from the Supabase
		top, err := u.calculatePriority(ctx, upc)
		if err != nil {
			slog.Error("Failed to calculate priority",
				slog.Group("upcomingTarget", "error", err),
			)
			return nil, err
		}
		video = top
	}

	upcChats, err := u.chatSvc.FetchChatsByVideoInfo(ctx, video, 0)
	if err != nil {
		slog.Error("Failed to fetch chats from the upcoming target video",
			slog.Group("upcomingTarget", "error", err),
		)
		return nil, err
	}

	// Save the fetched history to the Supabase
	if err := u.supaRepo.InsertFetchedHistory(ctx, video.SourceID); err != nil {
		slog.Error("Failed to insert fetched history",
			slog.Group("upcomingTarget", "sourceId", video.SourceID, "error", err),
		)
		return nil, err
	}

	// Filter chats by author channel
	targetChats, _ := filterChatsByAuthorChannel(upcChats, u.targetChannel)

	if len(targetChats) == 0 {
		slog.Info("No new chats from the upcoming target video",
			slog.Group("upcomingTarget", "sourceId", video.SourceID),
		)
		return nil, nil
	}

	// Filter chats by the publishedAt
	newChats, err := u.filterChatsByPublishedAt(ctx, targetChats, video.SourceID)
	if err != nil {
		slog.Error("Failed to filter chats by the publishedAt",
			slog.Group("upcomingTarget", "error", err),
		)
		return nil, err
	}

	return newChats, nil
}

func filterChatsByAuthorChannel(chats []model.YTChat, targetChannel []string) ([]model.YTChat, []model.YTChat) {
	var targetChats, otherChats []model.YTChat
	for _, chat := range chats {
		if slices.Contains(targetChannel, chat.AuthorChannelID) {
			targetChats = append(targetChats, chat)
		} else {
			otherChats = append(otherChats, chat)
		}
	}

	return targetChats, otherChats
}

func (u *chatUsecase) filterChatsByPublishedAt(ctx context.Context, chats []model.YTChat, sourceId string) ([]model.YTChat, error) {
	// Fetch the last recorded chat's publishedAt from the Supabase
	threshold, err := u.supaRepo.GetPublishedAtOfLastRecordedChatBySource(ctx, sourceId)
	if err != nil {
		slog.Error("Failed to get the last recorded chat",
			slog.Group("filterChat", "chats", chats, "sourceId", sourceId),
		)
		return nil, err
	}

	// Filter the chats by the threshold
	// The chats are already sorted by the publishedAt in ascending order (constraint of the YouTube API)

	var filteredChats []model.YTChat

	for i, chat := range chats {
		if chat.PublishedAtUnix > threshold {
			filteredChats = chats[i:]
			break
		}
	}

	return filteredChats, nil
}

func (u *chatUsecase) calculatePriority(ctx context.Context, videos []model.VideoInfo) (model.VideoInfo, error) {
	// Calculate priority
	// The priority is based on the fetched history

	ids := make([]string, 0, len(videos))
	for _, v := range videos {
		ids = append(ids, v.SourceID)
	}

	histories, err := u.supaRepo.GetFetchedHistory(ctx, ids)
	if err != nil {
		slog.Error("Failed to get fetched history",
			slog.Group("upcomingTarget", "error", err),
		)
		return model.VideoInfo{}, err
	}

	if len(histories) == 0 {
		// If all videos do not have fetched history, set the first video as the target
		return videos[0], nil
	}

	if len(histories) < len(videos) {
		// If there are videos with no history retrieved, set the first video with no history as the target
		for _, v := range videos {
			for _, h := range histories {
				if v.SourceID == h.SourceID {
					continue
				}
			}
			return v, nil
		}
	}

	// If all videos have fetched history, calculate the priority
	// The video with the oldest fetched history is the target

	oldest := histories[0]
	for _, v := range videos {
		if v.SourceID == oldest.SourceID {
			return v, nil
		}
	}

	slog.Error(
		"Failed to calculate priority",
		slog.Group("upcomingTarget",
			"histories", histories,
			"videos", videos,
			"error", fmt.Errorf("failed to calculate priority"),
		),
	)
	return model.VideoInfo{}, fmt.Errorf("failed to calculate priority")
}
