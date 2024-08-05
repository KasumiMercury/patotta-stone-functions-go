package realtime

import (
	"context"
	"fmt"
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/domain/video"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/status"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/test/migrate"
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

	if err := migrate.Migrate(connStr, "../../../../migrations"); err != nil {
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
			sourceIDs: []string{"get_source_id"},
			want: []*Record{
				{
					SourceID:    "get_source_id",
					Title:       "get_title",
					Status:      "archived",
					ChatID:      "get_chat_id",
					ScheduledAt: nil,
					UpdatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:      "empty",
			sourceIDs: []string{"non_existent_source_id"},
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
					"new_channel_id",
					"new_source_id",
					"new_title",
					"new_description",
					"new_chat_id",
					status.Archived,
					synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					synchro.In[tz.AsiaTokyo](time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
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
