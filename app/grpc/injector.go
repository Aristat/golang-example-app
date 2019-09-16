// +build wireinject

package grpc

import (
	"github.com/aristat/golang-example-app/app/config"
	"github.com/aristat/golang-example-app/app/entrypoint"
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/app/tracing"
	"github.com/google/wire"
)

// BuildPool
func Build() (*PoolManager, func(), error) {
	panic(wire.Build(ProviderProductionSet, entrypoint.ProviderProductionSet, logger.ProviderProductionSet, config.ProviderSet, tracing.ProviderProductionSet))
}

func BuildTest() (*PoolManager, func(), error) {
	panic(wire.Build(ProviderTestSet, entrypoint.ProviderTestSet, logger.ProviderTestSet, tracing.ProviderTestSet))
}
