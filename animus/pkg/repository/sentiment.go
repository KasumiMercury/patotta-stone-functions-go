package repository

import (
	"context"
)

type Sentiment interface {
	AnalyzeSentiment(ctx context.Context, text string) (float32, float32, error)
}
