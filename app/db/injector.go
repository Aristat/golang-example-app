// +build wireinject

package db

import (
	"github.com/aristat/golang-oauth2-example-app/app/provider"
	"github.com/google/wire"
)

func Build() (*Manager, func(), error) {
	panic(wire.Build(ProviderProductionSet, provider.AwareProductionSet))
}

func BuildTest() (*Manager, func(), error) {
	panic(wire.Build(ProviderTestSet, provider.AwareTestSet))
}
