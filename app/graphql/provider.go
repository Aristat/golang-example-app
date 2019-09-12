package graphql

import (
	"context"

	"github.com/aristat/golang-example-app/app/resolver"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/generated/graphql"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("graphql", &c)
	if e != nil {
		return c, func() {}, e
	}
	if cfg.IsSet("debug") && !c.Debug {
		c.Debug = cfg.GetBool("debug")
	}
	return c, func() {}, nil
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

// Provider
func Provider(ctx context.Context, resolver graphql.Config, log logger.Logger, cfg Config) (*GraphQL, func(), error) {
	g := New(ctx, resolver, log, cfg)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg, resolver.ProviderProductionSet)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest, resolver.ProviderTestSet)
)
