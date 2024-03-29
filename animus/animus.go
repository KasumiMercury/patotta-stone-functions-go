package animus

import (
	language "cloud.google.com/go/language/apiv2"
	"context"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/infra"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/lib"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/service"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/usecase"
	"github.com/uptrace/bun"
	"google.golang.org/api/youtube/v3"
	"log"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var stampPat *regexp.Regexp

// Global variables
// Initialize once per function instance
var ytApiKey = os.Getenv("YOUTUBE_API_KEY")
var ytSvc *youtube.Service

// DSN is the connection string for Supabase
var dsn = os.Getenv("SUPABASE_DSN")
var supaClient *bun.DB

var nlaClient *language.Client

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error

	// Define the pattern for the stamp
	// The pattern is `:stamp:`
	stampPat = regexp.MustCompile(`:[^:]+:`)

	// Custom log
	handler := lib.NewCustomLogger()
	slog.SetDefault(handler)

	// Create YouTube Data API service
	if ytApiKey == "" {
		slog.Error("YOUTUBE_API_KEY is not set")
		log.Fatalf("YOUTUBE_API_KEY is not set")
	}
	ytSvc, err = infra.NewYouTubeService(context.Background(), ytApiKey)
	if err != nil {
		slog.Error("Failed to create YouTube service", slog.Group("YouTubeAPI", "error", err))
		log.Fatalf("Failed to create YouTube service: %v", err)
	}

	// Create connection to Supabase
	if dsn == "" {
		slog.Error("DSN is not set")
		log.Fatalf("DSN is not set")
	}
	supaClient, err = infra.NewSupabaseClient(dsn)
	if err != nil {
		slog.Error("Failed to create Supabase client", slog.Group("Supabase", "error", err))
		log.Fatalf("Failed to create Supabase client: %v", err)
	}

	// Create connection to NaturalLanguageAPI
	nlaClient, err = infra.NewAnalysisClient(context.Background())
	if err != nil {
		slog.Error("Failed to create NaturalLanguageAPI client", slog.Group("NaturalLanguageAPI", "error", err))
		log.Fatalf("Failed to create NaturalLanguageAPI client: %v", err)
	}

	// Register the function to handle HTTP requests
	functions.HTTP("Animus", animus)
}

func animus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	ytRepo := infra.NewYouTubeRepository(ytSvc)
	supaRepo := infra.NewSupabaseRepository(supaClient)
	sntRepo := infra.NewSentimentRepository(nlaClient)

	chatSvc := service.NewChatService(ytRepo, supaRepo, sntRepo, stampPat)
	chatUsc := usecase.NewChatUsecase(targetChannels, chatSvc, supaRepo)

	videoUsc := usecase.NewVideoUsecase(supaRepo)

	// Get variable video info from Supabase
	targetStatus := []string{"upcoming", "live"}
	varVideos, err := videoUsc.GetVideoInfosByStatusFromSupabase(ctx, targetStatus)
	if err != nil {
		slog.Error("Failed to get variable video info", slog.Group("error", err))
		http.Error(w, "Failed to get video info by status", http.StatusInternalServerError)
		return
	}

	// Check the existence of the live status video
	// If there is live status video, skip the function
	if _, ok := varVideos["live"]; ok {
		slog.Info("There is a live status video")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Find the upcoming target video
	var upcVideos []model.VideoInfo
	if _, ok := varVideos["upcoming"]; ok {
		upcVideos = varVideos["upcoming"]
	} else {
		slog.Info("No upcoming target video")
	}

	// Fetch chats from the static target video
	stcChats, err := chatUsc.FetchChatsFromStaticTargetVideo(ctx)
	if err != nil {
		slog.Error("Failed to fetch chats from the static target video", slog.Group("staticTarget", "error", err))
		http.Error(w, "Failed to fetch chats from the static target video", http.StatusInternalServerError)
		return
	}

	// Fetch chats from the upcoming target video
	upcChats, err := chatUsc.FetchChatsFromUpcomingTargetVideo(ctx, upcVideos)
	if err != nil {
		slog.Error("Failed to fetch chats from the upcoming target video", slog.Group("upcomingTarget", "error", err))
		http.Error(w, "Failed to fetch chats from the upcoming target video", http.StatusInternalServerError)
		return
	}

	newChats := append(stcChats, upcChats...)
	if len(newChats) == 0 {
		slog.Info("No new chats")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Save the new target chats
	if err := chatUsc.SaveNewChats(ctx, newChats); err != nil {
		slog.Error("Failed to save new target chats",
			slog.Group("saveNewChats", "error", err),
		)
		http.Error(w, "Failed to save new target chats", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Animus function executed successfully")
}
