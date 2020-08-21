// +build wireinject

package graphql_resolver

import (
	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/aristat/golang-example-app/generated/graphql"
	"github.com/google/wire"
)

// Build
func Build() (graphql.Config, func(), error) {
	panic(wire.Build(ProviderProductionSet, repo.ProviderProductionSet, db.ProviderProductionSet, provider.AwareProductionSet))
}

func BuildTest() (graphql.Config, func(), error) {
	panic(wire.Build(ProviderTestSet, repo.ProviderTestSet, db.ProviderTestSet, provider.AwareTestSet))
}
