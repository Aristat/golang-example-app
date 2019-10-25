package products_router

import (
	"context"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/google/wire"
)

var ProviderManagers = wire.NewSet(
	wire.Struct(new(ServiceManagers), "*"),
)

// Provider
func Provider(ctx context.Context, log logger.Logger, managers ServiceManagers) (*Manager, func(), error) {
	g := New(ctx, log, managers)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderManagers)
	ProviderTestSet       = wire.NewSet(Provider, ProviderManagers)
)
