// +build wireinject

package users

import (
	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/oauth"
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/aristat/golang-example-app/app/session"
	"github.com/google/wire"
)

// Build
func Build() (*Manager, func(), error) {
	panic(wire.Build(ProviderProductionSet, oauth.ProviderProductionSet, session.ProviderProductionSet, repo.ProviderProductionSet, db.ProviderProductionSet, provider.AwareProductionSet))
}

func BuildTest() (*Manager, func(), error) {
	panic(wire.Build(ProviderTestSet, oauth.ProviderTestSet, session.ProviderTestSet, repo.ProviderTestSet, db.ProviderTestSet, provider.AwareTestSet))
}
