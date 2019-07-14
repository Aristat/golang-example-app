package session

import (
	"context"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"github.com/go-session/redis"
	"github.com/go-session/session"
)

// Config
type Config struct {
	RedisUrl string
	RedisDB  int
}

// New
func New(ctx context.Context, log *logger.Zap, cfg Config) (*session.Manager, func(), error) {
	log.Info("Initialize session manager")

	sessionStore := redis.NewRedisStore(&redis.Options{
		Addr: cfg.RedisUrl,
		DB:   cfg.RedisDB,
	})

	if _, err := sessionStore.Check(ctx, ""); err != nil {
		return nil, func() {}, err
	}

	sessionManager := session.NewManager(
		session.SetStore(sessionStore),
	)

	return sessionManager, func() {}, nil
}
