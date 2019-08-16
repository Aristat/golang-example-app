package users

import (
	"context"

	"github.com/aristat/golang-oauth2-example-app/app/db/domain"

	"github.com/aristat/golang-oauth2-example-app/app/db/repo"

	"github.com/aristat/golang-oauth2-example-app/app/db"
	"github.com/aristat/golang-oauth2-example-app/app/oauth"

	"github.com/aristat/golang-oauth2-example-app/app/logger"
	"github.com/go-session/session"
	"github.com/google/wire"
)

type Repo struct {
	Users domain.UsersRepo
}

var ProviderRepo = wire.NewSet(
	repo.NewAuthorsRepo,
	wire.Struct(new(Repo), "*"),
)

var ProviderTestRepo = wire.NewSet(
	repo.NewAuthorsRepo,
	wire.Struct(new(Repo), "*"),
)

// Managers
type Managers struct {
	session *session.Manager
	db      *db.Manager
	oauth   *oauth.Manager
}

var ProviderManagers = wire.NewSet(
	wire.Struct(new(Managers), "*"),
)

// Provider
func Provider(ctx context.Context, log logger.Logger, managers Managers, repo *Repo) (*Manager, func(), error) {
	g := New(ctx, log, managers, repo)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderRepo, ProviderManagers)
	ProviderTestSet       = wire.NewSet(Provider, ProviderTestRepo, ProviderManagers)
)
