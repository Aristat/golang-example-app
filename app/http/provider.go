package http

import (
	"context"
	"html/template"

	"github.com/aristat/golang-gin-oauth2-example-app/users"

	"github.com/aristat/golang-gin-oauth2-example-app/app/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/go-session/session"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"
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
func Mux(oauth *oauth.OAuth, db *db.Manager, session *session.Manager, log logger.Logger) (*chi.Mux, func(), error) {
	if mux != nil {
		return mux, func() {}, nil
	}

	mux = chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(Logger(log))

	tmp := template.Must(template.New("").ParseGlob("templates/**/*"))

	usersService := &users.Service{
		SessionManager: session,
		DB:             db.DB,
		Server:         oauth.OauthService.Server,
		Log:            log,
		Template:       tmp,
	}

	users.Run(mux, usersService)
	oauth.OauthService.Run(mux)

	return mux, func() {}, nil
}

// Provider
func Provider(ctx context.Context, mux *chi.Mux, log logger.Logger, cfg Config, oauth *oauth.OAuth, session *session.Manager) (*Http, func(), error) {
	g := New(ctx, mux, log, cfg, oauth, session)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg, Mux)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest)
)
