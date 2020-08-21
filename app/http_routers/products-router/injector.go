// +build wireinject

package products_router

import (
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/google/wire"
)

// Build
func Build() (*Manager, func(), error) {
	panic(wire.Build(ProviderProductionSet, provider.AwareProductionSet))
}

func BuildTest() (*Manager, func(), error) {
	panic(wire.Build(ProviderTestSet, provider.AwareTestSet))
}
