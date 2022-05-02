package db

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"

	"github.com/aristat/golang-example-app/app/logger"

	"github.com/google/wire"
	mocket "github.com/selvatico/go-mocket"
	"github.com/spf13/viper"
)

// Cfg
func Cfg(cfg *viper.Viper) (Config, func(), error) {
	c := Config{LogLevel: logger.LevelDebug}
	e := cfg.UnmarshalKey("db", &c)
	if e != nil {
		return c, func() {}, e
	}
	return c, func() {}, nil
}

// CfgTest
func CfgTest() (Config, func(), error) {
	return Config{}, func() {}, nil
}

// ProviderGORM
func ProviderGORM(ctx context.Context, log logger.Logger, cfg Config) (*gorm.DB, func(), error) {
	log = log.WithFields(logger.Fields{"service": prefix})

	db, err := gorm.Open("postgres", cfg.URL)
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns)
	db.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	db.DB().SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if cfg.LogLevel == logger.LevelDebug {
		db.LogMode(true)
	}

	cleanup := func() {
		if db != nil {
			_ = db.Close()
		}
		log.Info("Closed db connections")
	}

	return db, cleanup, err
}

// ProviderGORMTest for test, using go-mocket
func ProviderGORMTest() (*gorm.DB, func(), error) {
	var db *gorm.DB

	cleanup := func() {
		if db != nil {
			_ = db.Close()
		}
	}

	mocket.Catcher.Register()
	mocket.Catcher.Logging = true

	sqlDB, err := sql.Open(mocket.DriverName, "gorm")
	if err != nil {
		return db, cleanup, err
	}

	db, err = gorm.Open("postgres", sqlDB)

	return db, cleanup, err
}

// Provider
func Provider(ctx context.Context, log logger.Logger, cfg Config, db *gorm.DB) (*Manager, func(), error) {
	g := New(ctx, log, cfg, db)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderGORM, Cfg)
	ProviderTestSet       = wire.NewSet(Provider, ProviderGORMTest, CfgTest)
)
