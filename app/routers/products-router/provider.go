package products_router

import (
	"context"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Cfg
func Cfg(cfg *viper.Viper) (*Config, func(), error) {
	c := &Config{}
	e := cfg.UnmarshalKey("products", c)
	if e != nil {
		return c, func() {}, e
	}
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	return &Config{}, func() {}, nil
}

// Config
type Config struct {
	NatsURL string
	Subject string
}

var ProviderManagers = wire.NewSet(
	wire.Struct(new(ServiceManagers), "*"),
)

// Provider
func Provider(ctx context.Context, log logger.Logger, managers ServiceManagers, cfg *Config) (*Manager, func(), error) {
	g := New(ctx, log, managers, cfg)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderManagers, Cfg)
	ProviderTestSet       = wire.NewSet(Provider, ProviderManagers, CfgTest)
)
