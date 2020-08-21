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

// Provider
func Provider(ctx context.Context, cfg Config) (*session.Manager, func(), error) {
	return new(ctx, cfg)
}

func ProviderTest() (*session.Manager, func(), error) {
	return session.NewManager(), func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg)
	ProviderTestSet       = wire.NewSet(ProviderTest)
)
