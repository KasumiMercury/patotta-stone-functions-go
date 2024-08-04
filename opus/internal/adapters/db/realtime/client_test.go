package realtime

import (
	"context"
	"fmt"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/domain/video"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/test/testcontainers"
	"log"
	"os"
	"testing"
	"time"
)

var clt *Realtime

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := testcontainers.SetUpRealtimeContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// add SSL mode
	connStr += "sslmode=disable"
	fmt.Println(connStr)

	clt, err = NewRealtimeClient(connStr)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestRealtime_UpsertRecords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		video []video.Video
	}{
		{
			name: "success",
			video: func() []video.Video {
				v, _ := video.NewVideo(
					"channelID",
					"sourceID",
					"title",
					"description",
					"chatID",
					status.Archived,
					synchro.In[tz.AsiaTokyo](time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					synchro.In[tz.AsiaTokyo](time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					synchro.In[tz.AsiaTokyo](time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				)

				return []video.Video{*v}
			}(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := clt.UpsertRecords(context.Background(), tt.video); err != nil {
				t.Errorf("error: %v", err)
			}
		})
	}
}
