package grpc

import (
	"context"
	"errors"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/app/tracing"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const prefix = "app.grpc"

var errCfgInvalid = errors.New("cfg is not present or invalid")

// PoolManager
type PoolManager struct {
	ctx     context.Context
	tracing tracing.Tracer
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
	opts = append(opts, grpc.WithInsecure())

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
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			logger.UnaryClientInterceptor(p.logger, true),
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(p.tracing)),
		)))
	opts = append(opts, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
		logger.StreamClientInterceptor(p.logger, true),
		grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(p.tracing)),
	)))

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
func NewPoolManager(ctx context.Context, tracing tracing.Tracer, logger logger.Logger, cfg *Config) *PoolManager {
	return &PoolManager{ctx: ctx, tracing: tracing, cfg: cfg, logger: logger}
}
