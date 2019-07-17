package logger

import (
	"context"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// ProviderCfg returns configuration for production logger
func ProviderCfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("logger", &c)
	if e != nil {
		return c, func() {}, e
	}
	if cfg.IsSet("debug") {
		c.Debug = cfg.GetBool("debug")
	}
	return c, func() {}, nil
}

// ProviderCfgTest returns configuration for stub/mock logger
func ProviderCfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

// Provider returns logger instance implemented of Logger interface with resolved dependencies
func Provider(ctx context.Context, cfg Config) (*Zap, func(), error) {
	return NewZap(ctx, cfg), func() {}, nil
}

// ProviderTest returns stub/mock logger instance implemented of Logger interface with resolved dependencies
func ProviderTest(ctx context.Context, cfg Config) (*Mock, func(), error) {
	return NewMock(ctx, cfg, true), func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderCfg, wire.Bind(new(Logger), new(*Zap)))
	ProviderTestSet       = wire.NewSet(ProviderTest, ProviderCfgTest, wire.Bind(new(Logger), new(*Mock)))
)
