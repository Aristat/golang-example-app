package users

import (
	"context"
	"errors"
	"html/template"

	"github.com/aristat/golang-oauth2-example-app/app/logger"
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

func New(ctx context.Context, log logger.Logger, managers Managers, repo *Repo) *Manager {
	tmp := template.Must(template.New("").ParseGlob("templates/**/*"))
	log = log.WithFields(logger.Fields{"service": prefix})

	router := &Router{
		ctx:            ctx,
		sessionManager: managers.session,
		template:       tmp,
		logger:         log,
		db:             managers.db.DB,
		server:         managers.oauth.Router.Server,
		repo:           repo,
	}

	return &Manager{
		ctx:    ctx,
		logger: log,
		Router: router,
	}
}
