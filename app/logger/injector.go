// +build wireinject

package logger

import (
	"github.com/aristat/golang-example-app/app/config"
	"github.com/aristat/golang-example-app/app/entrypoint"
	"github.com/google/wire"
)

// Build returns logger instance implemented of Logger interface with resolved dependencies
func Build() (Logger, func(), error) {
	panic(wire.Build(ProviderProductionSet, entrypoint.ProviderProductionSet, config.ProviderSet))
}

func BuildTest() (Logger, func(), error) {
	panic(wire.Build(ProviderTestSet, entrypoint.ProviderTestSet))
}
