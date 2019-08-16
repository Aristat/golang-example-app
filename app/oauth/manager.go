package oauth

import (
	"context"

	"gopkg.in/oauth2.v3/store"

	"github.com/aristat/golang-oauth2-example-app/app/logger"
	"github.com/go-session/session"
	"gopkg.in/oauth2.v3"
)

const prefix = "app.oauth"

// OAuth
type Manager struct {
	ctx    context.Context
	Logger logger.Logger
	Router *Router
}

// New
func New(ctx context.Context, log logger.Logger, tokenStore oauth2.TokenStore, session *session.Manager, clientStore *store.ClientStore) *Manager {
	log = log.WithFields(logger.Fields{"service": prefix})

	server := NewServer(log, tokenStore, clientStore)
	router := &Router{
		ctx:            ctx,
		Server:         server,
		SessionManager: session,
	}
	router.Server.UserAuthorizationHandler = userAuthorization(router)

	return &Manager{
		ctx:    ctx,
		Logger: log,
		Router: router,
	}
}
