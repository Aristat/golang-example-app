// +build wireinject

package session

import (
	"github.com/aristat/golang-gin-oauth2-example-app/app/config"
	"github.com/aristat/golang-gin-oauth2-example-app/app/entrypoint"
	"github.com/go-session/session"
	"github.com/google/wire"
)

// Build
func Build() (*session.Manager, func(), error) {
	panic(wire.Build(ProviderProductionSet, config.ProviderSet, entrypoint.ProviderProductionSet))
}

// BuildTest
func BuildTest() (*session.Manager, func(), error) {
	panic(wire.Build(ProviderTestSet, entrypoint.ProviderTestSet))
}
