package repository

import "context"

type Supabase interface {
	GetLastUpdatedUnixOfVideo(ctx context.Context) (int64, error)
}
