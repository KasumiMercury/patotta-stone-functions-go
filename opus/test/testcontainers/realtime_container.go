package testcontainers

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type RealtimeContainer struct {
	*postgres.PostgresContainer
}

func SetUpRealtimeContainer(ctx context.Context) (*RealtimeContainer, error) {
	timeout := 2 * time.Minute

	container, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		postgres.WithDatabase("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(timeout),
		),
	)

	if err != nil {
		return nil, err
	}

	// migrate
	if _, _, err := container.Exec(ctx, []string{"psql", "-U", "postgres", "-d", "test", "-c", "CREATE TABLE videos (id SERIAL PRIMARY KEY, title TEXT, url TEXT, source_id TEXT, chat_id TEXT, status TEXT, scheduled_at TIMESTAMP, created_at TIMESTAMP, updated_at TIMESTAMP)"}); err != nil {
		return nil, err
	}

	return &RealtimeContainer{container}, nil
}