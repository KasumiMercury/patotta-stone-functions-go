package patotta_stone_functions_go

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/infra"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/service"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/usecase"
	"log/slog"
	"net/http"
	"os"
)

func init() {
	functions.HTTP("Animus", animus)
}

func animus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Cache common environment variables
	// Because the function is supposed to run on CloudFunctions, it is necessary to read the environment variables here.
	// If the environment variable is not set, the function will panic.
	// (To prevent retries by CloudScheduler, the function should panic without returning error responses.)
	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	if ytApiKey == "" {
		slog.Error("YOUTUBE_API_KEY is not set")
		panic("YOUTUBE_API_KEY must be set")
	}

	// Create YouTube service
	ytRepo, err := infra.NewYouTubeRepository(ctx, ytApiKey)
	if err != nil {
		slog.Error("Failed to create YouTube service", slog.Group("YouTubeAPI", "error", err))
		http.Error(w, "Failed to create YouTube service", http.StatusInternalServerError)
		return
	}
	chatSvc := service.NewChatService(ytRepo)
	chatUsc := usecase.NewChatUsecase(chatSvc)

	// Fetch chats from the static target video
	_, err = chatUsc.FetchChatsFromStaticTargetVideo(ctx)
	if err != nil {
		slog.Error("Failed to fetch chats from the static target video", slog.Group("staticTarget", "error", err))
		http.Error(w, "Failed to fetch chats from the static target video", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Animus function executed successfully")
}
