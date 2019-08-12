package session

import (
	"context"

	"github.com/go-session/redis"
	"github.com/go-session/session"
)

// Config
type Config struct {
	RedisUrl string
	RedisDB  int
}

// New
func New(ctx context.Context, cfg Config) (*session.Manager, func(), error) {
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
