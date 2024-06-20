package opus

import (
	"context"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/infra"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/lib"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/pkg/usecase"
	"github.com/uptrace/bun"
	"google.golang.org/api/youtube/v3"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// Global variables
// Initialize once per function instance
var ytApiKey = os.Getenv("YOUTUBE_API_KEY")
var ytSvc *youtube.Service

// DSN is the connection string for Supabase
var dsn = os.Getenv("SUPABASE_DSN")
var supaClient *bun.DB

func init() {
	// err is pre-declared to avoid a shadowing client.
	var err error

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
	ytRepo := infra.NewYouTubeRepository(ytSvc)
	rssRepo := infra.NewRssRepository()

	rssUsc := usecase.NewRssUsecase(rssRepo, supaRepo)
	apiSvc := usecase.NewApiUsecase(ytRepo)
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
		nsIds := make([]string, 0, len(pcs.NewItems))
		for _, p := range pcs.NewItems {
			nsIds = append(nsIds, p.SourceID)
			n = append(n, model.VideoInfo{
				SourceID:      p.SourceID,
				Title:         p.Title,
				UpdatedAtUnix: p.UpdatedAtUnix,
			})
		}
		// Fetch details of new videos
		vdMap, err := apiSvc.VideoDetailsMap(ctx, nsIds)
		err = videoUsc.SaveNewVideo(ctx, n, vdMap)

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

	slog.Info("Fetched updated RSS items")

	// extract parameter from request
	q := r.URL.Query()
	// extract upcoming param
	uw := q.Get("upcoming")
	// if uw is not empty, check update scheduler of upcoming videos
	if uw != "" {
		slog.Info("Check update scheduler of upcoming videos")
		// extract upcoming status video
		uv := make([]model.VideoSchedule, 0, len(recMap))
		for _, record := range recMap {
			if record.Status == "upcoming" {
				uv = append(uv, model.VideoSchedule{
					SourceID:    record.SourceID,
					Status:      record.Status,
					ScheduledAt: record.ScheduledAt,
				})
			}
		}
		// if there is no upcoming status video, return
		if len(uv) == 0 {
			slog.Info("No upcoming status video")
			w.WriteHeader(http.StatusOK)
			return
		}

		// fetch YouTube api and compare scheduled time
		uIds := make([]string, 0, len(uv))
		for _, u := range uv {
			uIds = append(uIds, u.SourceID)
		}
		// fetch scheduled time of upcoming videos
		vSaMap, err := apiSvc.VideoScheduledAtUnixMap(ctx, uIds)
		if err != nil {
			slog.Error("Failed to fetch scheduled time of upcoming videos",
				slog.Group("check upcoming schedule", "error", err),
			)
			http.Error(w, "Failed to fetch scheduled time of upcoming videos", http.StatusInternalServerError)
			return
		}

		// compare scheduled time
		// if scheduled time is different, update scheduled time
		// if scheduled time is same, do nothing
		// if scheduled time is not found, update status to archived
		for _, u := range uv {
			if sa, ok := vSaMap[u.SourceID]; ok {
				if u.ScheduledAt.Unix() != sa {
					err := videoUsc.UpdateScheduledAt(ctx, u.SourceID, time.Unix(sa, 0))
					if err != nil {
						slog.Error("Failed to update scheduled time",
							slog.Group("check upcoming schedule", "error", err),
						)
						http.Error(w, "Failed to update scheduled time", http.StatusInternalServerError)
						return
					}
				}
			} else {
				// TODO: update status to archived
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
