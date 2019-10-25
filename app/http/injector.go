// +build wireinject

package http

import (
	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/graphql"
	"github.com/aristat/golang-example-app/app/oauth"
	"github.com/aristat/golang-example-app/app/provider"
	users_router "github.com/aristat/golang-example-app/app/routers/users-router"
	"github.com/aristat/golang-example-app/app/session"
	"github.com/google/wire"
)

// Build
func Build() (*Http, func(), error) {
	panic(wire.Build(
		ProviderProductionSet,
		graphql.ProviderProductionSet,
		users_router.ProviderProductionSet,
		oauth.ProviderProductionSet,
		session.ProviderProductionSet,
		repo.ProviderProductionSet,
		db.ProviderProductionSet,
		provider.AwareProductionSet,
	))
}
