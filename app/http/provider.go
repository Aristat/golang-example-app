package http

import (
	"context"

	"github.com/riandyrn/otelchi"

	"github.com/aristat/golang-example-app/app/dataloader"

	products_router "github.com/aristat/golang-example-app/app/http_routers/products-router"

	"github.com/aristat/golang-example-app/app/auth"

	"github.com/aristat/golang-example-app/app/graphql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aristat/golang-example-app/app/logger"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var muxRouter *chi.Mux

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("http", &c)
	if e != nil {
		return c, func() {}, e
	}
	c.Debug = cfg.GetBool("debug")
	return c, func() {}, nil
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

// Mux
func Mux(managers Managers, log logger.Logger) (*chi.Mux, func(), error) {
	if muxRouter != nil {
		return muxRouter, func() {}, nil
	}

	muxRouter = chi.NewRouter()
	muxRouter.Use(middleware.RequestID)
	muxRouter.Use(Logger(log))
	muxRouter.Use(otelchi.Middleware("http-server", otelchi.WithChiRoutes(muxRouter)))
	muxRouter.Use(dataloader.LoaderMiddleware)

	managers.products.Router.Run(muxRouter)
	managers.graphql.Routers(muxRouter.With(managers.authMiddleware.JWTHandler))

	return muxRouter, func() {}, nil
}

// ServiceManagers
type Managers struct {
	products       *products_router.Manager
	authMiddleware *auth.Middleware
	graphql        *graphql.GraphQL
}

var ProviderManagers = wire.NewSet(
	wire.Struct(new(Managers), "*"),
)

// Provider
func Provider(ctx context.Context, mux *chi.Mux, log logger.Logger, cfg Config) (*Http, func(), error) {
	g := New(ctx, mux, log, cfg)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg, Mux, ProviderManagers)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest)
)
