package db

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"
	_ "github.com/lib/pq"
)

const prefix = "app.db"

// Config
type Config struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	LogLevel        logger.Level
}

// Http
type Manager struct {
	ctx context.Context
	cfg Config
	log logger.Logger
	DB  *gorm.DB
}

// New
func New(ctx context.Context, log logger.Logger, cfg Config, db *gorm.DB) *Manager {
	return &Manager{
		ctx: ctx,
		cfg: cfg,
		log: log.WithFields(logger.Fields{"service": prefix}),
		DB:  db,
	}
}
