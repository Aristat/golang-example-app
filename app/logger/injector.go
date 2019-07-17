// +build wireinject

package logger

import (
	"github.com/aristat/golang-gin-oauth2-example-app/app/config"
	"github.com/aristat/golang-gin-oauth2-example-app/app/entrypoint"
	"github.com/google/wire"
)

// Build returns logger instance implemented of Logger interface with resolved dependencies
func Build() (Logger, func(), error) {
	panic(wire.Build(ProviderProductionSet, entrypoint.ProviderProductionSet, config.ProviderSet))
}
