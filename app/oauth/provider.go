package oauth

import (
	"context"

	"gopkg.in/oauth2.v3/store"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"github.com/go-session/session"

	"github.com/google/wire"
	"github.com/spf13/viper"

	"gopkg.in/oauth2.v3"

	oauthRedis "gopkg.in/go-oauth2/redis.v3"
)

// Config
type Config struct {
	RedisUrl string
	RedisDB  int
}

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("oauth", &c)
	if e != nil {
		return c, func() {}, nil
	}
	return c, func() {}, nil
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

func TokenStore(cfg Config) (oauth2.TokenStore, func(), error) {
	oauthConfig := oauthRedis.Options{
		Addr: cfg.RedisUrl,
		DB:   cfg.RedisDB,
	}

	tokenStore := oauthRedis.NewRedisStore(&oauthConfig)
	return tokenStore, func() {}, nil
}

func TokenStoreTest() (oauth2.TokenStore, func(), error) {
	tokenStore, err := store.NewMemoryTokenStore()
	return tokenStore, func() {}, err
}

// Provider
func Provider(ctx context.Context, log logger.Logger, tokenStore oauth2.TokenStore, session *session.Manager) (*Manager, func(), error) {
	g := New(ctx, log, tokenStore, session)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg, TokenStore)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest, TokenStoreTest)
)
