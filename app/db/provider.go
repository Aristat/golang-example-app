package db

import (
	"context"
	"database/sql"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{}
	e := cfg.UnmarshalKey("db", &c)
	if e != nil {
		return c, func() {}, nil
	}
	return c, func() {}, nil
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

func DB(cfg Config, log *logger.Zap) (*sql.DB, func(), error) {
	log.Info("Initialize DB")

	db, err := sql.Open("postgres", cfg.URL)

	if err != nil {
		return nil, func() {}, err
	}

	if err := db.Ping(); err != nil {
		return nil, func() {}, err
	}

	return db, func() {}, nil
}

// Provider
func Provider(ctx context.Context, log *logger.Zap, cfg Config, db *sql.DB) (*Manager, func(), error) {
	g := New(ctx, log, cfg, db)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, DB, Cfg)
	ProviderTestSet       = wire.NewSet(Provider, DB, CfgTest)
)
