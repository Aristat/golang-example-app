// +build wireinject

package graphql

import (
	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/google/wire"
)

// Build
func Build() (*GraphQL, func(), error) {
	panic(wire.Build(ProviderProductionSet, repo.ProviderProductionSet, db.ProviderProductionSet, provider.AwareProductionSet))
}

// BuildTest
func BuildTest() (*GraphQL, func(), error) {
	panic(wire.Build(ProviderTestSet, repo.ProviderTestSet, db.ProviderTestSet, provider.AwareTestSet))
}
