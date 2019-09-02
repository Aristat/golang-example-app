package oauth

import (
	"context"

	"github.com/go-oauth2/oauth2/models"

	"gopkg.in/oauth2.v3/store"

	"github.com/aristat/golang-example-app/app/logger"

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

func ClientStore() (*store.ClientStore, func(), error) {
	clientsConfig := map[string]oauth2.ClientInfo{
		"123456": &models.Client{
			ID:     "123456",
			Secret: "12345678",
			Domain: "http://localhost:9094",
		},
	}

	clientStore, err := NewClientStore(clientsConfig)
	return clientStore, func() {}, err
}

func ClientStoreTest() (*store.ClientStore, func(), error) {
	clientsConfig := map[string]oauth2.ClientInfo{
		"123456": &models.Client{
			ID:     "123456",
			Secret: "12345678",
			Domain: "http://127.0.0.1:8090",
		},
	}

	clientStore, err := NewClientStore(clientsConfig)
	return clientStore, func() {}, err
}

func NewClientStore(config map[string]oauth2.ClientInfo) (*store.ClientStore, error) {
	clientStore := store.NewClientStore()
	for key, value := range config {
		err := clientStore.Set(key, value)

		if err != nil {
			return nil, err
		}
	}

	return clientStore, nil
}

// Provider
func Provider(ctx context.Context, log logger.Logger, tokenStore oauth2.TokenStore, session *session.Manager, clientStore *store.ClientStore) (*Manager, func(), error) {
	g := New(ctx, log, tokenStore, session, clientStore)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg, TokenStore, ClientStore)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest, TokenStoreTest, ClientStoreTest)
)
