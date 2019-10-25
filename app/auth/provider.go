package auth

import (
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// ProviderCfg
func ProviderCfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("auth", &c)
	return c, func() {}, e
}

// Provider
func Provider(cfg Config, logger logger.Logger) (*Middleware, func(), error) {
	return NewMiddleware(cfg, logger)
}

// ProviderTest
func ProviderTest(logger logger.Logger) (*Middleware, func(), error) {
	return NewTestMiddleware(logger)
}

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderCfg)
	ProviderTestSet       = wire.NewSet(ProviderTest)
)
