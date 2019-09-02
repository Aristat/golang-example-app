// +build wireinject

package http

import (
	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/graphql"
	"github.com/aristat/golang-example-app/app/oauth"
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/aristat/golang-example-app/app/session"
	"github.com/aristat/golang-example-app/app/users"
	"github.com/google/wire"
)

// Build
func Build() (*Http, func(), error) {
	panic(wire.Build(
		ProviderProductionSet,
		graphql.ProviderProductionSet,
		users.ProviderProductionSet,
		oauth.ProviderProductionSet,
		session.ProviderProductionSet,
		repo.ProviderProductionSet,
		db.ProviderProductionSet,
		provider.AwareProductionSet,
	))
}
