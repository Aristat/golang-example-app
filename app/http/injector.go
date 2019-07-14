// +build wireinject

package http

import (
	"github.com/aristat/golang-gin-oauth2-example-app/app/config"
	"github.com/aristat/golang-gin-oauth2-example-app/app/entrypoint"
	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"
	"github.com/aristat/golang-gin-oauth2-example-app/app/session"
	"github.com/google/wire"
)

// Build
func Build() (*Http, func(), error) {
	panic(wire.Build(ProviderProductionSet, oauth.ProviderProductionSet, session.ProviderProductionSet, config.ProviderSet, entrypoint.ProviderProductionSet))
}
