package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/service"
	"log/slog"
	"os"
)

type Chat interface {
	FetchChatsFromStaticTargetVideo(ctx context.Context) ([]model.YTChat, error)
}

type chatUsecase struct {
	chatSvc service.Chat
}

func NewChatUsecase(chatSvc service.Chat) Chat {
	return &chatUsecase{
		chatSvc: chatSvc,
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

	// debug log
	slog.Info("Fetched chats from the static target video", "count", len(stcChats))

	return stcChats, nil
}
