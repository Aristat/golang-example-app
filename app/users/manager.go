package users

import (
	"context"
	"errors"
	"html/template"

	"github.com/aristat/golang-gin-oauth2-example-app/app/db"

	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"
	"github.com/go-session/session"
)

var (
	userNotFound = errors.New("10002 user not found")
)

type H map[string]interface{}

const prefix = "app.users"

// OAuth
type Manager struct {
	ctx    context.Context
	logger logger.Logger
	Router *Router
}

func New(ctx context.Context, log logger.Logger, db *db.Manager, session *session.Manager, oauth *oauth.Manager) *Manager {
	tmp := template.Must(template.New("").ParseGlob("templates/**/*"))
	log = log.WithFields(logger.Fields{"service": prefix})

	router := &Router{
		ctx:            ctx,
		sessionManager: session,
		template:       tmp,
		logger:         log,
		db:             db.DB,
		server:         oauth.Router.Server,
	}

	return &Manager{
		ctx:    ctx,
		logger: log,
		Router: router,
	}
}
