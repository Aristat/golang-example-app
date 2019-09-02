package resolver

import (
	"context"

	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/generated/graphql"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("resolver", &c)
	return c, func() {}, e
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

type Config struct {
}

// Managers
type Managers struct {
	Repo *repo.Repo
}

var ProviderManagers = wire.NewSet(
	wire.Struct(new(Managers), "*"),
)

// Provider
func Provider(ctx context.Context, log logger.Logger, cfg Config, managers Managers) (graphql.Config, func(), error) {
	c := New(ctx, log, cfg, managers)
	return c, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg, ProviderManagers)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest, ProviderManagers)
)
