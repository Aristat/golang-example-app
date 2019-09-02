// +build wireinject

package session

import (
	"github.com/aristat/golang-example-app/app/provider"
	"github.com/go-session/session"
	"github.com/google/wire"
)

// Build
func Build() (*session.Manager, func(), error) {
	panic(wire.Build(ProviderProductionSet, provider.AwareProductionSet))
}

func BuildTest() (*session.Manager, func(), error) {
	panic(wire.Build(ProviderTestSet))
}
