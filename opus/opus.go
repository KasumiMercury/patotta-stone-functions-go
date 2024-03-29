package opus

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/infra"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/lib"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/usecase"
	"github.com/uptrace/bun"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// DSN is the connection string for Supabase
var dsn = os.Getenv("SUPABASE_DSN")
var supaClient *bun.DB

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error

	// Custom log
	handler := lib.NewCustomLogger()
	slog.SetDefault(handler)

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

	supaRepo := infra.NewSupabaseRepository(supaClient)

	rssRepo := infra.NewRssRepository()
	rssUsc := usecase.NewRssUsecase(rssRepo, supaRepo)

	videoUsc := usecase.NewVideoUsecase(supaRepo)

	rss, err := rssUsc.FetchUpdatedRssItemsEachOfChannels(ctx, targetChannels)
	if err != nil {
		slog.Error("Failed to fetch updated RSS items", slog.Group("rssWatch", "error", err))
		http.Error(w, "Failed to fetch updated RSS items", http.StatusInternalServerError)
		return
	}
	slog.Debug("Fetched updated RSS items", slog.Group("rssWatch", "rss", rss))

	if len(rss) == 0 {
		slog.Info("No updated RSS items")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Store sourceIDs of videos that have been updated in an array
	sIds := make([]string, 0, len(rss))
	for _, r := range rss {
		sIds = append(sIds, r.SourceID)
	}

	// Get existing video records
	recMap, err := videoUsc.GetVideoRecordMap(ctx, sIds)
	if err != nil {
		slog.Error("Failed to get video records", slog.Group("rssWatch", "error", err))
		http.Error(w, "Failed to get video records", http.StatusInternalServerError)
		return
	}

	// Compare RSS items and video records
	pcs := rssUsc.CompareRssItemsAndVideoRecords(ctx, rss, recMap)

	// Save new videos
	if len(pcs.NewItems) > 0 {
		n := make([]model.VideoInfo, 0, len(pcs.NewItems))
		for _, p := range pcs.NewItems {
			n = append(n, model.VideoInfo{
				SourceID:      p.SourceID,
				Title:         p.Title,
				UpdatedAtUnix: p.UpdatedAtUnix,
			})
		}
		err = videoUsc.SaveNewVideo(ctx, n)

		if err != nil {
			slog.Error("Failed to save new videos",
				slog.Group("rssWatch",
					slog.Group("saveNewVideo", "error", err),
				),
			)
			http.Error(w, "Failed to save new videos", http.StatusInternalServerError)
			return
		}
		slog.Debug(
			"Saved new videos",
			slog.Group("rssWatch", "newVideos", n),
		)
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Fetched updated RSS items")
}
