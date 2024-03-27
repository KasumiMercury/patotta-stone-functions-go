package infra

import (
	language "cloud.google.com/go/language/apiv2"
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

func NewSentimentRepository(ctx context.Context) (*SentimentRepository, error) {
	client, err := NewAnalysisClient(ctx)
	if err != nil {
		slog.Error(
			"Failed to create a new language client",
			slog.Group("NaturalLanguageAPI", "error", err),
		)
		return nil, err
	}
	return &SentimentRepository{client: client}, nil
}
