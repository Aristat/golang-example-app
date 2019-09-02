package http

import (
	"context"

	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/graphql"
	"github.com/aristat/golang-example-app/app/users"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/go-session/session"

	"github.com/aristat/golang-example-app/app/logger"

	"github.com/aristat/golang-example-app/app/oauth"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var mux *chi.Mux

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("http", &c)
	if e != nil {
		return c, func() {}, nil
	}
	c.Debug = cfg.GetBool("debug")
	return c, func() {}, nil
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

// Mux
func Mux(db *db.Manager, managers Managers, log logger.Logger) (*chi.Mux, func(), error) {
	if mux != nil {
		return mux, func() {}, nil
	}

	mux = chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(Logger(log))

	managers.users.Router.Run(mux)
	managers.oauth.Router.Run(mux)
	managers.graphql.Routers(mux)

	return mux, func() {}, nil
}

// Managers
type Managers struct {
	session *session.Manager
	users   *users.Manager
	oauth   *oauth.Manager
	graphql *graphql.GraphQL
}

var ProviderManagers = wire.NewSet(
	wire.Struct(new(Managers), "*"),
)

// Provider
func Provider(ctx context.Context, mux *chi.Mux, log logger.Logger, cfg Config, managers Managers) (*Http, func(), error) {
	g := New(ctx, mux, log, cfg, managers)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg, Mux, ProviderManagers)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest)
)
