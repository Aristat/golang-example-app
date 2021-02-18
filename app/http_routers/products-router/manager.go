package products_router

import (
	"context"

	"github.com/aristat/golang-example-app/app/grpc"
	"github.com/aristat/golang-example-app/app/logger"
)

const prefix = "app.products-router"

// Product Manager
type Manager struct {
	ctx    context.Context
	logger logger.Logger
	Router *Router
}

// ServiceManagers
type ServiceManagers struct {
	PoolManager *grpc.PoolManager
}

func New(ctx context.Context, log logger.Logger, managers ServiceManagers, cfg *Config) *Manager {
	log = log.WithFields(logger.Fields{"service": prefix})

	router := &Router{
		ctx:         ctx,
		cfg:         cfg,
		logger:      log,
		poolManager: managers.PoolManager,
	}

	return &Manager{
		ctx:    ctx,
		logger: log,
		Router: router,
	}
}
