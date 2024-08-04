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

	// migrate
	if _, _, err := container.Exec(ctx, []string{"psql", "-U", "postgres", "-d", "test", "-c", "CREATE TABLE videos (id SERIAL PRIMARY KEY, title TEXT, url TEXT, source_id TEXT, chat_id TEXT, status TEXT, scheduled_at TIMESTAMP, created_at TIMESTAMP, updated_at TIMESTAMP)"}); err != nil {
		log.Fatal(err)
	}
	// unique index
	if _, _, err := container.Exec(ctx, []string{"psql", "-U", "postgres", "-d", "test", "-c", "CREATE UNIQUE INDEX idx_videos_source_id ON videos (source_id)"}); err != nil {
		log.Fatal(err)
	}

	// insert test data
	if _, _, err := container.Exec(ctx, []string{"psql", "-U", "postgres", "-d", "test", "-c", "INSERT INTO videos (title, url, source_id, chat_id, status, scheduled_at, created_at, updated_at) VALUES ('title', 'url', 'sourceID', 'chatID', 'status', '2022-01-01 00:00:00', '2022-01-01 00:00:00', '2022-01-01 00:00:00')"}); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestRealtime_GetRecordsBySourceIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		sourceIDs []string
		want      []*Record
	}{
		{
			name:      "success",
			sourceIDs: []string{"sourceID"},
			want: []*Record{
				{
					SourceID:    "sourceID",
					Title:       "title",
					Status:      "archived",
					ChatID:      "chatID",
					ScheduledAt: timeToPtr(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
					UpdatedAt:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:      "empty",
			sourceIDs: []string{"sourceID2"},
			want:      nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := clt.GetRecordsBySourceIDs(context.Background(), tt.sourceIDs)
			if err != nil {
				t.Errorf("error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
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

func timeToPtr(t time.Time) *time.Time {
	return &t
}
