// +build wireinject

package auth

import (
	"github.com/aristat/golang-example-app/app/config"
	"github.com/google/wire"
)

// Build
func Build() (*Middleware, func(), error) {
	panic(wire.Build(ProviderProductionSet, config.ProviderSet))
}

// BuildTest
func BuildTest() (*Middleware, func(), error) {
	panic(wire.Build(ProviderTestSet))
}
