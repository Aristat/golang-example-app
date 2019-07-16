package session

import (
	"context"

	"github.com/go-session/session"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("session", &c)
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
func Provider(ctx context.Context, cfg Config) (*session.Manager, func(), error) {
	return New(ctx, cfg)
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest)
)
