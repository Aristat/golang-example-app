// +build wireinject

package oauth_router

import (
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/aristat/golang-example-app/app/session"
	"github.com/google/wire"
)

// Build
func Build() (*Manager, func(), error) {
	panic(wire.Build(ProviderProductionSet, session.ProviderProductionSet, provider.AwareProductionSet))
}

func BuildTest() (*Manager, func(), error) {
	panic(wire.Build(ProviderTestSet, session.ProviderTestSet, provider.AwareTestSet))
}
