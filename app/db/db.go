package db

import (
	"context"
	"database/sql"
)

// Config
type Config struct {
	URL string
}

// Http
type DBManager struct {
	ctx context.Context
	cfg Config
	DB  *sql.DB
}

// New
func New(ctx context.Context, cfg Config, db *sql.DB) *DBManager {
	return &DBManager{
		ctx: ctx,
		cfg: cfg,
		DB:  db,
	}
}
