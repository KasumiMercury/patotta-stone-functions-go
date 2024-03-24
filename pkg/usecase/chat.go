package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/repository"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/service"
	"log/slog"
	"os"
	"slices"
)

type Chat interface {
	FetchChatsFromStaticTargetVideo(ctx context.Context) ([]model.YTChat, error)
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

	// Filter chats by the publishedAt
	newChats, err := u.filterChatsByPublishedAt(ctx, targetChats, stc.SourceID)
	if err != nil {
		slog.Error("Failed to filter chats by the publishedAt", "error", err)
		return nil, err
	}

	// Save the new chats to the Supabase
	if err := u.chatSvc.SaveNewTargetChats(ctx, newChats); err != nil {
		slog.Error("Failed to insert the new chats", "error", err)
		return nil, err
	}

	// debug log
	slog.Info("Fetched chats from the static target video", "count", len(newChats))

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
		slog.Error("Failed to get the last recorded chat", "error", err)
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
