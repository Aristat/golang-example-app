package users

import (
	"context"

	"github.com/aristat/golang-example-app/app/db/repo"

	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/oauth"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/go-session/session"
	"github.com/google/wire"
)

// Managers
type Managers struct {
	Session *session.Manager
	DB      *db.Manager
	Oauth   *oauth.Manager
	Repo    *repo.Repo
}

var ProviderManagers = wire.NewSet(
	wire.Struct(new(Managers), "*"),
)

// Provider
func Provider(ctx context.Context, log logger.Logger, managers Managers) (*Manager, func(), error) {
	g := New(ctx, log, managers)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderManagers)
	ProviderTestSet       = wire.NewSet(Provider, ProviderManagers)
)
