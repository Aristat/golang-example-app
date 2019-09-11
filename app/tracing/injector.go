// +build wireinject

package tracing

import (
	"github.com/aristat/golang-example-app/app/config"
	"github.com/aristat/golang-example-app/app/entrypoint"
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/google/wire"
)

// Build
func Build() (Tracer, func(), error) {
	panic(wire.Build(ProviderProductionSet, entrypoint.ProviderProductionSet, logger.ProviderProductionSet, config.ProviderSet))
}
