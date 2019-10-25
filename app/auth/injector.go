// +build wireinject

package auth

import (
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/google/wire"
)

// Build
func Build() (*Middleware, func(), error) {
	panic(wire.Build(ProviderProductionSet, provider.AwareProductionSet))
}

// BuildTest
func BuildTest() (*Middleware, func(), error) {
	panic(wire.Build(ProviderTestSet, provider.AwareTestSet))
}
