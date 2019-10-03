package auth

import (
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
func Provider(cfg Config) (*Middleware, func(), error) {
	return NewMiddleware(cfg)
}

// ProviderTest
func ProviderTest() (*Middleware, func(), error) {
	return NewTestMiddleware()
}

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderCfg)
	ProviderTestSet       = wire.NewSet(ProviderTest)
)
