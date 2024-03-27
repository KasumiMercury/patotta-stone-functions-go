package infra

import (
	language "cloud.google.com/go/language/apiv2"
	"cloud.google.com/go/language/apiv2/languagepb"
	"context"
	"log/slog"
)

type SentimentRepository struct {
	client *language.Client
}

func NewAnalysisClient(ctx context.Context) (*language.Client, error) {
	client, err := language.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewSentimentRepository(client *language.Client) *SentimentRepository {
	return &SentimentRepository{client: client}
}

func (r *SentimentRepository) AnalyzeSentiment(ctx context.Context, text string) (float32, float32, error) {
	sentiment, err := r.client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
	if err != nil {
		slog.Error(
			"Failed to analyze sentiment",
			slog.Group("analyzeSentiment", "text", text,
				slog.Group("NaturalLanguageAPI", "error", err),
			),
		)
		return 0, 0, err
	}
	return sentiment.DocumentSentiment.Score, sentiment.DocumentSentiment.Magnitude, nil
}
