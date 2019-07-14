// +build wireinject

package logger

import (
	"github.com/aristat/golang-gin-oauth2-example-app/app/config"
	"github.com/aristat/golang-gin-oauth2-example-app/app/entrypoint"
	"github.com/google/wire"
)

func Build() (*Zap, func(), error) {
	panic(wire.Build(ProviderProductionSet, entrypoint.ProviderProductionSet, config.ProviderSet))
}
