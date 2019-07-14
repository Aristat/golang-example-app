// +build wireinject

package oauth

import (
	"github.com/aristat/golang-gin-oauth2-example-app/app/config"
	"github.com/aristat/golang-gin-oauth2-example-app/app/entrypoint"
	"github.com/aristat/golang-gin-oauth2-example-app/app/session"
	"github.com/google/wire"
)

// Build
func Build() (*OAuth, func(), error) {
	panic(wire.Build(ProviderProductionSet, session.ProviderProductionSet, config.ProviderSet, entrypoint.ProviderProductionSet))
}
