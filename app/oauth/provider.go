package oauth

import (
	"context"

	"github.com/go-session/session"

	"github.com/google/wire"
	"github.com/spf13/viper"
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

// Provider
func Provider(ctx context.Context, cfg Config, session *session.Manager) (*OAuth, func(), error) {
	g := New(ctx, cfg, session)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg)
)
