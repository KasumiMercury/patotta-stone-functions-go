package service

import (
	"context"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	model2 "github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/model"
	repository2 "github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/repository"
	"log/slog"
	"time"
)

type Chat interface {
	FetchChatsByVideoInfo(ctx context.Context, videoInfo model2.VideoInfo, l int64) ([]model2.YTChat, error)
	SaveNewTargetChats(ctx context.Context, chats []model2.YTChat) error
}

type chatService struct {
	ytRepo   repository2.YouTube
	supaRepo repository2.Supabase
}

func NewChatService(ytRepo repository2.YouTube, supaRepo repository2.Supabase) Chat {
	return &chatService{
		ytRepo:   ytRepo,
		supaRepo: supaRepo,
	}
}

func (s *chatService) FetchChatsByVideoInfo(ctx context.Context, videoInfo model2.VideoInfo, l int64) ([]model2.YTChat, error) {
	// Fetch chats from the static target video
	resp, err := s.ytRepo.FetchChatsByChatID(ctx, videoInfo.ChatID, l)
	if err != nil {
		slog.Error(
			"Failed to fetch chats",
			slog.Group("fetchChat", "chatId", videoInfo.ChatID, "sourceId", videoInfo.SourceID,
				slog.Group("YouTubeAPI", "error", err),
			),
		)
		return nil, err
	}

	items := resp.Items
	if items == nil || len(items) == 0 {
		return nil, nil
	}

	// Create a slice to store the chats
	chats := make([]model2.YTChat, 0, len(resp.Items))
	for _, item := range items {
		pa, err := synchro.ParseISO[tz.AsiaTokyo](item.Snippet.PublishedAt)
		if err != nil {
			slog.Error(
				"Failed to parse the publishedAt",
				slog.Group("fetchChat", "chatId", videoInfo.ChatID, "sourceId", videoInfo.SourceID,
					slog.Group("formatting", "error", err, "target", item.Snippet.PublishedAt)),
			)
			continue
		}
		chats = append(chats, model2.YTChat{
			AuthorChannelID: item.Snippet.AuthorChannelId,
			Message:         item.Snippet.DisplayMessage,
			PublishedAtUnix: pa.Unix(),
			SourceID:        videoInfo.SourceID,
		})
	}

	return chats, nil
}

func (s *chatService) SaveNewTargetChats(ctx context.Context, chats []model2.YTChat) error {
	// Convert the chats to the chat records
	chatRecords := make([]model2.ChatRecord, 0, len(chats))
	for _, chat := range chats {
		chatRecords = append(chatRecords, model2.ChatRecord{
			Message:     chat.Message,
			SourceID:    chat.SourceID,
			PublishedAt: time.Unix(chat.PublishedAtUnix, 0),
		})
	}

	// Save the chats to the database
	if err := s.supaRepo.InsertChatRecord(ctx, chatRecords); err != nil {
		slog.Error("Failed to insert chats",
			slog.Group("saveChat", "sourceId", chats[0].SourceID,
				slog.Group("Supabase", "error", err),
			),
		)
		return err
	}

	return nil
}
