package common

import (
	"log"

	"context"

	"github.com/go-session/redis"
	"github.com/go-session/session"
	"github.com/spf13/viper"
)

func InitSession() {
	redisURL := viper.GetString("REDIS_URL")
	db := viper.GetInt("REDIS_SESSION_DB")

	log.Printf("[REDIS] Init redis: %v, db: %v", redisURL, db)

	manager := redis.NewRedisStore(&redis.Options{
		Addr: redisURL,
		DB:   int(db),
	})

	_, err := manager.Check(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}

	session.InitManager(
		session.SetStore(manager),
	)
}
