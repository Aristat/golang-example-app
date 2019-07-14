package entrypoint

import (
	"context"

	"github.com/google/wire"
)

// ContextProvider returns stub/mock context instance with resolved dependencies
func ContextProvider() (context.Context, func(), error) {
	c := OnShutdown()
	return c, func() {}, nil
}

// ContextProviderTest returns stub/mock context instance with resolved dependencies
func ContextProviderTest() (context.Context, func(), error) {
	c := context.Background()
	return c, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(ContextProvider)
	ProviderTestSet       = wire.NewSet(ContextProviderTest)
)
