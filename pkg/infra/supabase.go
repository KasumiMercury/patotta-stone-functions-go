package infra

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type SupabaseRepository struct {
	db *bun.DB
}

func NewSupabaseClient(dsn string) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return db, nil
}

func NewSupabaseRepository(dsn string) (*SupabaseRepository, error) {
	db, err := NewSupabaseClient(dsn)
	if err != nil {
		return nil, err
	}

	return &SupabaseRepository{db: db}, nil
}
