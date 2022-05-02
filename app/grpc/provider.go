package grpc

import (
	"context"
	"sync"
	"time"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/aristat/golang-example-app/app/logger"

	"github.com/google/wire"
	"github.com/spf13/viper"
	"google.golang.org/grpc/keepalive"
)

var (
	pm      *PoolManager
	mutexPM sync.Mutex
)

// Cfg
func Cfg(cfg *viper.Viper) (*Config, func(), error) {
	c := &Config{}
	e := cfg.UnmarshalKey("grpc", c)
	if e != nil {
		return c, func() {}, e
	}
	return c, func() {}, e
}

// CfgTest
func CfgTest() (*Config, func(), error) {
	return &Config{}, func() {}, nil
}

// Service
type Service struct {
	Target           string
	MaxConn          int
	InitConn         int
	MaxLifeDuration  time.Duration
	IdleTimeout      time.Duration
	ClientParameters *keepalive.ClientParameters
}

// Config
type Config struct {
	Services         map[string]*Service
	ClientParameters *keepalive.ClientParameters
}

// Provider
func Provider(ctx context.Context, tracing *tracesdk.TracerProvider, logger logger.Logger, cfg *Config) (*PoolManager, func(), error) {
	mutexPM.Lock()
	defer mutexPM.Unlock()
	if pm != nil {
		return pm, func() {}, nil
	}
	pm = NewPoolManager(ctx, tracing, logger, cfg)
	return pm, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest)
)
