package http

import (
	"context"

	"github.com/aristat/golang-gin-oauth2-example-app/app/db"

	"github.com/go-session/session"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

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

// Provider
func Provider(ctx context.Context, log *logger.Zap, cfg Config, oath *oauth.OAuth, session *session.Manager, db *db.Manager) (*Http, func(), error) {
	g := New(ctx, log, cfg, oath, session, db)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg)
	ProviderTestSet       = wire.NewSet(Provider, CfgTest)
)
