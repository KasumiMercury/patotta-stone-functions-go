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
	FetchChatsFromStaticTargetVideo(ctx context.Context) error
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

func (u *chatUsecase) FetchChatsFromStaticTargetVideo(ctx context.Context) error {
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
		return err
	}

	// Filter chats by author channel
	targetChats, _ := filterChatsByAuthorChannel(stcChats, u.targetChannel)

	// Filter chats by the publishedAt
	newChats, err := u.filterChatsByPublishedAt(ctx, targetChats, stc.SourceID)
	if err != nil {
		slog.Error("Failed to filter chats by the publishedAt",
			slog.Group("staticTarget", "error", err),
		)
		return err
	}

	if len(newChats) == 0 {
		slog.Info("No new chats from the static target video")
		return nil
	}

	// Save the new chats to the Supabase
	if err := u.chatSvc.SaveNewTargetChats(ctx, newChats); err != nil {
		slog.Error("Failed to insert the new chats",
			slog.Group("staticTarget",
				slog.Group("saveChat", "sourceId", newChats[0].SourceID, "error", err),
			),
		)
		return err
	}

	// debug log
	slog.Debug("Fetched chats from the static target video", "count", len(newChats))

	return nil
}

func (u *chatUsecase) FetchChatsFromUpcomingTargetVideo(ctx context.Context) error {
	// Get the upcoming target video info from the Supabase
	upc, err := u.supaRepo.GetVideoInfoByStatus(ctx, []string{"upcoming"})
	if err != nil {
		slog.Error("Failed to get the upcoming target video",
			slog.Group("upcomingTarget", "error", err),
		)
		return err
	}

	if len(upc) == 0 {
		slog.Info("No upcoming target video")
		return nil
	}

	// Fetch chats from the upcoming target video
	for _, video := range upc {
		info := model.VideoInfo{
			SourceID: video.SourceID,
			ChatID:   video.ChatID,
		}
		upcChats, err := u.chatSvc.FetchChatsByVideoInfo(ctx, info, 0)
		if err != nil {
			slog.Error("Failed to fetch chats from the upcoming target video",
				slog.Group("upcomingTarget", "error", err),
			)
			return err
		}

		// Filter chats by author channel
		targetChats, _ := filterChatsByAuthorChannel(upcChats, u.targetChannel)

		// Filter chats by the publishedAt
		newChats, err := u.filterChatsByPublishedAt(ctx, targetChats, video.SourceID)
		if err != nil {
			slog.Error("Failed to filter chats by the publishedAt",
				slog.Group("upcomingTarget", "error", err),
			)
			return err
		}

		if len(newChats) == 0 {
			slog.Info("No new chats from the upcoming target video")
			continue
		}

		// Save the new chats to the Supabase
		if err := u.chatSvc.SaveNewTargetChats(ctx, newChats); err != nil {
			slog.Error("Failed to insert the new chats",
				slog.Group("upcomingTarget",
					slog.Group("saveChat", "sourceId", newChats[0].SourceID, "error", err),
				),
			)
			return err
		}

		// debug log
		slog.Debug("Fetched chats from the upcoming target video", "count", len(newChats))
	}

	return nil
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
