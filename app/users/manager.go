package users

import (
	"context"
	"errors"
	"html/template"

	"github.com/aristat/golang-example-app/app/logger"
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

func New(ctx context.Context, log logger.Logger, managers Managers) *Manager {
	tmp := template.Must(template.New("").ParseGlob("templates/**/*"))
	log = log.WithFields(logger.Fields{"service": prefix})

	router := &Router{
		ctx:            ctx,
		sessionManager: managers.Session,
		template:       tmp,
		logger:         log,
		db:             managers.DB.DB,
		server:         managers.Oauth.Router.Server,
		repo:           managers.Repo,
		poolManager:    managers.PoolManager,
	}

	return &Manager{
		ctx:    ctx,
		logger: log,
		Router: router,
	}
}
