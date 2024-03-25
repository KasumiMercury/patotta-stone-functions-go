package animus

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/infra"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/lib"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/service"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/usecase"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func init() {
	functions.HTTP("Animus", animus)
}

func animus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Custom log
	handler := lib.NewCustomLogger()
	slog.SetDefault(handler)

	// Cache common environment variables
	// Because the function is supposed to run on CloudFunctions, it is necessary to read the environment variables here.
	// If the environment variable is not set, the function will panic.
	// (To prevent retries by CloudScheduler, the function should panic without returning error responses.)
	targetChannelIdStr := os.Getenv("TARGET_CHANNEL_ID")
	if targetChannelIdStr == "" {
		slog.Error("TARGET_CHANNEL_ID is not set")
		panic("TARGET_CHANNEL_ID is not set")
	}
	// Split targetChannelIdStr by comma
	targetChannels := strings.Split(targetChannelIdStr, ",")

	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	if ytApiKey == "" {
		slog.Error("YOUTUBE_API_KEY is not set")
		panic("YOUTUBE_API_KEY must be set")
	}

	dsn := os.Getenv("SUPABASE_DSN")
	if dsn == "" {
		slog.Error("DSN is not set")
		panic("DSN is not set")
	}

	// Create YouTube service
	ytRepo, err := infra.NewYouTubeRepository(ctx, ytApiKey)
	if err != nil {
		slog.Error("Failed to create YouTube service", slog.Group("YouTubeAPI", "error", err))
		http.Error(w, "Failed to create YouTube service", http.StatusInternalServerError)
		return
	}
	// Create Supabase client
	supaRepo, err := infra.NewSupabaseRepository(dsn)
	if err != nil {
		slog.Error("Failed to create Supabase client", slog.Group("Supabase", "error", err))
		http.Error(w, "Failed to create Supabase client", http.StatusInternalServerError)
		return
	}

	chatSvc := service.NewChatService(ytRepo, supaRepo)
	chatUsc := usecase.NewChatUsecase(targetChannels, chatSvc, supaRepo)

	videoUsc := usecase.NewVideoUsecase(supaRepo)

	// Get variable video info from Supabase
	targetStatus := []string{"upcoming"}
	varVideos, err := videoUsc.GetVideoInfosByStatusFromSupabase(ctx, targetStatus)
	if err != nil {
		slog.Error("Failed to get variable video info", slog.Group("error", err))
		http.Error(w, "Failed to get video info by status", http.StatusInternalServerError)
	}

	// Find the upcoming target video
	var upcVideos []model.VideoInfo
	if _, ok := varVideos["upcoming"]; ok {
		upcVideos = varVideos["upcoming"]
	} else {
		slog.Info("No upcoming target video")
	}

	// Fetch chats from the static target video
	err = chatUsc.FetchChatsFromStaticTargetVideo(ctx)
	if err != nil {
		slog.Error("Failed to fetch chats from the static target video", slog.Group("staticTarget", "error", err))
		http.Error(w, "Failed to fetch chats from the static target video", http.StatusInternalServerError)
		return
	}

	// Fetch chats from the upcoming target video
	err = chatUsc.FetchChatsFromUpcomingTargetVideo(ctx)
	if err != nil {
		slog.Error("Failed to fetch chats from the upcoming target video", slog.Group("upcomingTarget", "error", err))
		http.Error(w, "Failed to fetch chats from the upcoming target video", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Animus function executed successfully")
}
