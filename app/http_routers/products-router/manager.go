package products_router

import (
	"context"
	"html/template"

	"github.com/aristat/golang-example-app/app/entrypoint"
	"github.com/aristat/golang-example-app/app/grpc"
	"github.com/aristat/golang-example-app/app/logger"
)

const prefix = "app.products-router"

// OAuth Manager
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
	wd := entrypoint.WorkDir()
	tmp := template.Must(template.New("").ParseGlob(wd + "/templates/**/*"))
	log = log.WithFields(logger.Fields{"service": prefix})

	router := &Router{
		ctx:         ctx,
		cfg:         cfg,
		template:    tmp,
		logger:      log,
		poolManager: managers.PoolManager,
	}

	return &Manager{
		ctx:    ctx,
		logger: log,
		Router: router,
	}
}
