package graphql_resolver

import (
	"context"

	"github.com/casbin/casbin"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/generated/graphql"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("graphql_resolver", &c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{ProductTimeout: 5}, func() {}, nil
}

var ProviderManagers = wire.NewSet(
	wire.Struct(new(Managers), "*"),
)

// Provider
func Provider(ctx context.Context, log logger.Logger, cfg Config, enforcer *casbin.Enforcer, managers Managers) (graphql.Config, func(), error) {
	c := New(ctx, log, cfg, enforcer, managers)
	return c, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg, ProviderManagers)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest, ProviderManagers)
)
