package users

import (
	"context"

	"github.com/aristat/golang-gin-oauth2-example-app/app/db"
	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"
	"github.com/go-session/session"
	"github.com/google/wire"
)

// Provider
func Provider(ctx context.Context, log logger.Logger, db *db.Manager, session *session.Manager, oauth *oauth.Manager) (*Manager, func(), error) {
	g := New(ctx, log, db, session, oauth)
	return g, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider)
	ProviderTestSet       = wire.NewSet(Provider)
)
