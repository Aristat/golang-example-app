package grpc

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/aristat/golang-example-app/app/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const prefix = "app.grpc"

var errCfgInvalid = errors.New("cfg is not present or invalid")

// PoolManager
type PoolManager struct {
	ctx     context.Context
	tracing *tracesdk.TracerProvider
	cfg     *Config
	logger  logger.Logger
}

// New
func (p *PoolManager) NewPool(service string) (_ *Pool, loaded bool, _ error) {
	s, ok := p.cfg.Services[service]

	if !ok {
		return nil, false, errCfgInvalid
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	cl := s.ClientParameters
	if cl == nil {
		if p.cfg.ClientParameters != nil {
			cl = p.cfg.ClientParameters
		} else {
			cl = &keepalive.ClientParameters{}
		}
	}

	opts = append(opts, grpc.WithKeepaliveParams(*cl))
	opts = append(opts,
		grpc.WithChainUnaryInterceptor(
			logger.UnaryClientInterceptor(p.logger, true),
			otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(p.tracing), otelgrpc.WithPropagators(propagation.TraceContext{})),
		),
	)
	opts = append(opts,
		grpc.WithChainStreamInterceptor(
			logger.StreamClientInterceptor(p.logger, true),
			otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(p.tracing), otelgrpc.WithPropagators(propagation.TraceContext{})),
		),
	)

	pool, l := NewPool(p.ctx, service, s.Target,
		MaxConn(s.MaxConn),
		InitConn(s.InitConn),
		MaxLifeDuration(s.MaxLifeDuration),
		IdleTimeout(s.IdleTimeout),
		ConnOptions(opts...),
	)

	return pool, l, nil
}

// NewPoolManager
func NewPoolManager(ctx context.Context, tracing *tracesdk.TracerProvider, logger logger.Logger, cfg *Config) *PoolManager {
	return &PoolManager{ctx: ctx, tracing: tracing, cfg: cfg, logger: logger}
}
