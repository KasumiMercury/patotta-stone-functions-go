package opus

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/infra"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/usecase"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func init() {
	// Register the function to handle HTTP requests
	functions.HTTP("Opus", opus)
}

func opus(w http.ResponseWriter, r *http.Request) {
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

	rssRepo := infra.NewRssRepository()
	rssUsecase := usecase.NewRssUsecase(rssRepo)
	if err := rssUsecase.FetchUpdatedRssItemsEachOfChannels(ctx, targetChannels); err != nil {
		slog.Error("Failed to fetch updated RSS items", slog.Group("rssWatch", "error", err))
		http.Error(w, "Failed to fetch updated RSS items", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Fetched updated RSS items")
}
