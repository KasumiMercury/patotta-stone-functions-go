package service

import (
	"context"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/lib"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/model"
	"github.com/KasumiMercury/patotta-stone-functions-go/animus/pkg/repository"
	"golang.org/x/text/unicode/norm"
	"log/slog"
	"regexp"
	"time"
)

type Chat interface {
	FetchChatsByVideoInfo(ctx context.Context, videoInfo model.VideoInfo, l int64) ([]model.YTChat, error)
	SaveNewTargetChats(ctx context.Context, chats []model.YTChat) error
}

type chatService struct {
	ytRepo   repository.YouTube
	supaRepo repository.Supabase
	sntRepo  repository.Sentiment
	stampPat *regexp.Regexp
}

func NewChatService(ytRepo repository.YouTube, supaRepo repository.Supabase, sntRepo repository.Sentiment, stampPat *regexp.Regexp) Chat {
	return &chatService{
		ytRepo:   ytRepo,
		supaRepo: supaRepo,
		sntRepo:  sntRepo,
		stampPat: stampPat,
	}
}

func (s *chatService) FetchChatsByVideoInfo(ctx context.Context, videoInfo model.VideoInfo, l int64) ([]model.YTChat, error) {
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
	chats := make([]model.YTChat, 0, len(resp.Items))
	for _, item := range items {
		pa, err := synchro.ParseISO[tz.AsiaTokyo](item.Snippet.PublishedAt)
		if err != nil {
			slog.Error(
				"Failed to parse the publishedAt",
				slog.Group("fetchChat", "chatId", videoInfo.ChatID, "sourceId", videoInfo.SourceID,
					slog.Group("formatChat", "error", err, "target", item.Snippet.PublishedAt)),
			)
			continue
		}
		chats = append(chats, model.YTChat{
			AuthorChannelID: item.Snippet.AuthorChannelId,
			Message:         item.Snippet.DisplayMessage,
			PublishedAtUnix: pa.Unix(),
			SourceID:        videoInfo.SourceID,
		})
	}

	return chats, nil
}

func (s *chatService) SaveNewTargetChats(ctx context.Context, chats []model.YTChat) error {
	// Convert the chats to the chat records
	chatRecords := make([]model.ChatRecord, 0, len(chats))
	for _, chat := range chats {
		isNegative, err := s.analyzeNegativityOfChatMessage(ctx, chat.Message)
		if err != nil {
			slog.Error("Failed to analyze negativity of chat message",
				slog.Group("saveChat", "sourceId", chat.SourceID, "message", chat.Message, "error", err),
			)
			continue
		}
		chatRecords = append(chatRecords, model.ChatRecord{
			Message:     chat.Message,
			IsNegative:  isNegative,
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

func (s *chatService) analyzeNegativityOfChatMessage(ctx context.Context, msg string) (bool, error) {
	// Remove stamps from the chat message
	// Stamps are not necessary for sentiment analysis
	rsMsg := s.stampPat.ReplaceAllString(msg, "")
	// Normalize the message
	nMsg := norm.NFKC.String(rsMsg)
	// Remove emojis from the message
	// Because the emojis are not necessary for the sentiment analysis and occasionally cause an error
	reMsg := lib.RemoveEmoji(nMsg)

	score, magnitude, err := s.sntRepo.AnalyzeSentiment(ctx, reMsg)
	if err != nil {
		slog.Error("Failed to analyze sentiment",
			slog.Group("analyzeSentiment", "error", err),
		)
		return false, err
	}

	// Calculate the negativity of the chat message
	// Threshold is proportional to magnitude
	// The higher the magnitude, the more likely it is to contain negative elements, so the higher the threshold is set.
	negativity := score < 0.5*magnitude

	return negativity, nil
}
