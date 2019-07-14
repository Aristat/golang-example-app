package logger

import (
	"context"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Config is a general logger config settings
type Config struct {
	Debug   bool
	Verbose bool
}

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
	if cfg.IsSet("verbose") {
		c.Verbose = cfg.GetBool("verbose")
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

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderCfg)
)
