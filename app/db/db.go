package db

import (
	"context"
	"database/sql"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"
)

const prefix = "app.db"

// Config
type Config struct {
	URL string
}

// Http
type Manager struct {
	ctx context.Context
	cfg Config
	log *logger.Zap
	DB  *sql.DB
}

// New
func New(ctx context.Context, log *logger.Zap, cfg Config, db *sql.DB) *Manager {
	return &Manager{
		ctx: ctx,
		cfg: cfg,
		log: log.WithFields(logger.Fields{"service": prefix}),
		DB:  db,
	}
}
