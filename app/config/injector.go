// +build wireinject

package config

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Build returns spf13/viper instance with resolved dependencies
func Build() (*viper.Viper, func(), error) {
	panic(wire.Build(Provider))
}
