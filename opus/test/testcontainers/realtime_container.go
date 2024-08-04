package testcontainers

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type realtimeContainer struct {
	testcontainers.Container
}

func SetUpRealtimeContainer(ctx context.Context) (*realtimeContainer, error) {
	port, _ := nat.NewPort("tcp", "8080")
	timeout := 2 * time.Minute
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{port.Port()},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForListeningPort(port).WithStartupTimeout(timeout),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	return &realtimeContainer{container}, nil
}
