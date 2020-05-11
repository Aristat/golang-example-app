package entrypoint

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/wire"
)

// ContextProvider
func ContextProvider() (context.Context, func(), error) {
	c := OnShutdown()
	return c, func() {}, nil
}

// ContextProviderTest
func ContextProviderTest() (context.Context, func(), error) {
	// initializing
	wd, _ := filepath.Abs(os.Getenv("APP_WD"))
	_, err := Initialize(wd, nil)

	c := context.Background()
	return c, func() {}, err
}

var (
	ProviderProductionSet = wire.NewSet(ContextProvider)
	ProviderTestSet       = wire.NewSet(ContextProviderTest)
)
