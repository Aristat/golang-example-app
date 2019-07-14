// +build wireinject

package http

import (
	"github.com/aristat/golang-gin-oauth2-example-app/app/db"
	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"
	"github.com/aristat/golang-gin-oauth2-example-app/app/provider"
	"github.com/aristat/golang-gin-oauth2-example-app/app/session"
	"github.com/google/wire"
)

// Build
func Build() (*Http, func(), error) {
	panic(wire.Build(ProviderProductionSet, oauth.ProviderProductionSet, session.ProviderProductionSet, db.ProviderProductionSet, provider.AwareProductionSet))
}
