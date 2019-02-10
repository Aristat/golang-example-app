package common

import (
	"go/build"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetDefault("GOPATH", build.Default.GOPATH)

	viper.SetDefault("HTTP_SERVER_PORT", DefaultHttpServerPort)

	viper.SetDefault("REDIS_URL", DefaultRedisUrl)
	viper.SetDefault("REDIS_SESSION_DB", DefaultRedisDB)
	viper.SetDefault("REDIS_TOKEN_DB", DefaultRedisDB)

	viper.SetDefault("DATABASE_URL", DefaultDatabaseUrl)
	viper.SetDefault("TEST_DATABASE_URL", DefaultTestDatabaseUrl)
}
