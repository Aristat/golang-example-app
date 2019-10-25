// +build wireinject

package http

import (
	"github.com/aristat/golang-example-app/app/auth"
	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/graphql"
	"github.com/aristat/golang-example-app/app/provider"
	oauth_router "github.com/aristat/golang-example-app/app/routers/oauth-router"
	products_router "github.com/aristat/golang-example-app/app/routers/products-router"
	users_router "github.com/aristat/golang-example-app/app/routers/users-router"
	"github.com/aristat/golang-example-app/app/session"
	"github.com/google/wire"
)

// Build
func Build() (*Http, func(), error) {
	panic(wire.Build(
		ProviderProductionSet,
		auth.ProviderProductionSet,
		graphql.ProviderProductionSet,
		users_router.ProviderProductionSet,
		products_router.ProviderProductionSet,
		oauth_router.ProviderProductionSet,
		session.ProviderProductionSet,
		repo.ProviderProductionSet,
		db.ProviderProductionSet,
		provider.AwareProductionSet,
	))
}
