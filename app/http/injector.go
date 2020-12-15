// +build wireinject

package http

import (
	"github.com/aristat/golang-example-app/app/auth"
	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/graphql"
	products_router "github.com/aristat/golang-example-app/app/http_routers/products-router"
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/google/wire"
)

// Build
func Build() (*Http, func(), error) {
	panic(wire.Build(
		ProviderProductionSet,
		auth.ProviderProductionSet,
		graphql.ProviderProductionSet,
		products_router.ProviderProductionSet,
		repo.ProviderProductionSet,
		db.ProviderProductionSet,
		provider.AwareProductionSet,
	))
}
