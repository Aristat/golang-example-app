package db

import (
	"context"
	"database/sql"

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

func DB(cfg Config) (*sql.DB, func(), error) {
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
func Provider(ctx context.Context, cfg Config, db *sql.DB) (*DBManager, func(), error) {
	g := New(ctx, cfg, db)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, DB, Cfg)
	ProviderTestSet       = wire.NewSet(Provider, DB, CfgTest)
)
