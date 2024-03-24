package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/service"
	"log/slog"
	"os"
	"slices"
)

type Chat interface {
	FetchChatsFromStaticTargetVideo(ctx context.Context) ([]model.YTChat, error)
}

type chatUsecase struct {
	chatSvc       service.Chat
	targetChannel []string
}

func NewChatUsecase(chatSvc service.Chat, targetChannel []string) Chat {
	return &chatUsecase{
		chatSvc:       chatSvc,
		targetChannel: targetChannel,
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

	// debug log
	slog.Info("Fetched chats from the static target video", "count", len(targetChats))

	return targetChats, nil
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
