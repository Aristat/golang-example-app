// +build wireinject

package db

import (
	"github.com/aristat/golang-gin-oauth2-example-app/app/config"
	"github.com/aristat/golang-gin-oauth2-example-app/app/entrypoint"
	"github.com/google/wire"
)

func Build() (*DBManager, func(), error) {
	panic(wire.Build(ProviderProductionSet, config.ProviderSet, entrypoint.ProviderProductionSet))
}
