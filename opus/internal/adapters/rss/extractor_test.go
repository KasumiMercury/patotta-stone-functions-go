package rss

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRssClient_FetchRssItem(t *testing.T) {
	// read test data
	data, err := os.ReadFile("../../../test/videos.xml")
	if err != nil {
		t.Fatal(err)
	}

	// mock Rss Server
	mockServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/rss+xml")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(data)
				if err != nil {
					t.Fatal(err)
				}
			}),
	)
	defer mockServer.Close()

	// test FetchRssItems

	// arrange
	client := NewRssClient(NewParser())
	ctx := context.Background()

	tests := map[string]struct {
		LimitUnix int64
		ExpectLen int
	}{
		"limitUnix is 0": {
			LimitUnix: 0,
			ExpectLen: 6,
		},
		"limitTime is now": {
			LimitUnix: time.Now().Unix(),
			ExpectLen: 0,
		},
		"limitTime is 2024-06-21": {
			LimitUnix: time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC).Unix(),
			ExpectLen: 4,
		},
	}

	for name, tt := range tests {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			// act
			items, err := client.FetchRssItems(ctx, mockServer.URL, tt.LimitUnix)

			// assert
			if err != nil {
				t.Errorf("error: %v", err)
			}

			if len(items) != tt.ExpectLen {
				t.Errorf("got: %d, want: %d", len(items), tt.ExpectLen)
			}
		})
	}
}
