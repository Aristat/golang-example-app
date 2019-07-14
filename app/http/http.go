package http

import (
	"context"
	"syscall"

	"github.com/go-session/session"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"

	"github.com/aristat/golang-gin-oauth2-example-app/app/db"

	"github.com/aristat/golang-gin-oauth2-example-app/users"
	"github.com/fvbock/endless"

	"github.com/gin-gonic/gin"
)

const prefix = "app.http"

// Config
type Config struct {
	Debug bool
	Bind  string
}

// Http
type Http struct {
	ctx     context.Context
	cfg     Config
	oauth   *oauth.OAuth
	session *session.Manager
	db      *db.Manager
	log     *logger.Zap
}

// ListenAndServe
func (m *Http) ListenAndServe(bind ...string) (err error) {
	m.log.Info("Initialize http")

	bindAdrr := m.cfg.Bind

	if len(bind) > 0 && len(bind[0]) > 0 {
		bindAdrr = bind[0]
	}

	usersService := &users.Service{
		SessionManager: m.session,
		DB:             m.db.DB,
		OauthServer:    m.oauth.OauthServer,
	}

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("templates/**/*")

	users.Run(r, usersService)
	oauth.Run(r, m.oauth.OauthService)

	m.log.Info("start listen and serve http at %v", logger.Args(bindAdrr))

	server := endless.NewServer(bindAdrr, r)
	server.BeforeBegin = func(add string) {
		m.log.Info("Actual pid is %d", logger.Args(syscall.Getpid()))
	}

	go func() {
		<-m.ctx.Done()
		m.log.Info("context cancelled, shutdown is raised")
		if e := server.Shutdown(context.Background()); e != nil {
			m.log.Emergency("graceful shutdown error, %v", logger.Args(e))
		}
	}()

	err = server.ListenAndServe()
	if err != nil {
		m.log.Emergency("Server err: %v", logger.Args(err))
	}

	return
}

// New
func New(ctx context.Context, log *logger.Zap, cfg Config, oauth *oauth.OAuth, session *session.Manager, db *db.Manager) *Http {
	return &Http{
		ctx:     ctx,
		cfg:     cfg,
		oauth:   oauth,
		session: session,
		db:      db,
		log:     log.WithFields(logger.Fields{"service": prefix}),
	}
}
