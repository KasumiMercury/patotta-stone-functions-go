package patotta_stone_functions_go

import (
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/pkg/usecase"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
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
	ytSvc, err := youtube.NewService(ctx, option.WithAPIKey(ytApiKey))
	if err != nil {
		slog.Error("Failed to create YouTube service", slog.Group("YouTubeAPI", "error", err))
		http.Error(w, "Failed to create YouTube service", http.StatusInternalServerError)
		return
	}

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
	stcChats, err := usecase.FetchChatsByVideoInfo(ctx, ytSvc, stc, 0)
	if err != nil {
		slog.Error("Failed to fetch chats from the static target video", slog.Group("staticTarget", "error", err))
		http.Error(w, "Failed to fetch chats from the static target video", http.StatusInternalServerError)
		return
	}

	// debug log
	fmt.Printf("Fetched %d chats from the static target video", len(stcChats))

	w.WriteHeader(http.StatusOK)
	slog.Info("Animus function executed successfully")
}
