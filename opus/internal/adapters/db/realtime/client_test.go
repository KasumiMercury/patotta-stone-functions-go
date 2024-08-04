package realtime

import (
	"context"
	"fmt"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/test/testcontainers"
	"log"
	"os"
	"testing"
)

var clt *Realtime

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := testcontainers.SetUpRealtimeContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	host, _ := container.Host(ctx)
	p, _ := container.MappedPort(ctx, "8080/tcp")

	connectionString := fmt.Sprintf("postgres://postgres:password@%s:%s/test?sslmode=disable", host, p)

	clt, err = NewRealtimeClient(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}
