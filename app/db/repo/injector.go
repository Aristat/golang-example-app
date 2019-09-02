// +build wireinject

package repo

import (
	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/google/wire"
)

func Build() (*Repo, func(), error) {
	panic(wire.Build(ProviderProductionSet, db.ProviderProductionSet, provider.AwareProductionSet))
}

func BuildTest() (*Repo, func(), error) {
	panic(wire.Build(ProviderTestSet, db.ProviderTestSet))
}
