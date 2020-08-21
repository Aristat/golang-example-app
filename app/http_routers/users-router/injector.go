// +build wireinject

package users_router

import (
	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/repo"
	oauth_router "github.com/aristat/golang-example-app/app/http_routers/oauth-router"
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/aristat/golang-example-app/app/session"
	"github.com/google/wire"
)

// Build
func Build() (*Manager, func(), error) {
	panic(wire.Build(ProviderProductionSet, oauth_router.ProviderProductionSet, session.ProviderProductionSet, repo.ProviderProductionSet, db.ProviderProductionSet, provider.AwareProductionSet))
}

func BuildTest() (*Manager, func(), error) {
	panic(wire.Build(ProviderTestSet, oauth_router.ProviderTestSet, session.ProviderTestSet, repo.ProviderTestSet, db.ProviderTestSet, provider.AwareTestSet))
}
