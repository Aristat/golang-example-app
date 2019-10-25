package users_router

import (
	"context"
	"errors"
	"html/template"

	"github.com/aristat/golang-example-app/app/db"
	"github.com/aristat/golang-example-app/app/db/repo"
	"github.com/aristat/golang-example-app/app/entrypoint"
	"github.com/aristat/golang-example-app/app/grpc"
	oauth_router "github.com/aristat/golang-example-app/app/routers/oauth-router"
	"github.com/go-session/session"

	"github.com/aristat/golang-example-app/app/logger"
)

var (
	userNotFound = errors.New("10002 user not found")
)

type H map[string]interface{}

const prefix = "app.users-router"

// OAuth Manager
type Manager struct {
	ctx    context.Context
	logger logger.Logger
	Router *Router
}

// ServiceManagers
type ServiceManagers struct {
	Session     *session.Manager
	DB          *db.Manager
	Oauth       *oauth_router.Manager
	Repo        *repo.Repo
	PoolManager *grpc.PoolManager
}

func New(ctx context.Context, log logger.Logger, managers ServiceManagers) *Manager {
	wd := entrypoint.WorkDir()
	tmp := template.Must(template.New("").ParseGlob(wd + "/templates/**/*"))
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
