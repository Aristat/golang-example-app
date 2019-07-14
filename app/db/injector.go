// +build wireinject

package db

import (
	"github.com/aristat/golang-gin-oauth2-example-app/app/provider"
	"github.com/google/wire"
)

func Build() (*Manager, func(), error) {
	panic(wire.Build(ProviderProductionSet, provider.AwareProductionSet))
}
