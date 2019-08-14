package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/go-session/session"

	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"

	"github.com/aristat/golang-gin-oauth2-example-app/app/oauth"
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
	oauth   *oauth.Manager
	session *session.Manager
	log     logger.Logger
	mux     *chi.Mux
}

// ListenAndServe
func (m *Http) ListenAndServe(bind ...string) (err error) {
	bindAddress := m.cfg.Bind

	if len(bind) > 0 && len(bind[0]) > 0 {
		bindAddress = bind[0]
	}

	server := &http.Server{
		Addr:    bindAddress,
		Handler: m.mux,
	}

	go func() {
		<-m.ctx.Done()
		m.log.Info("context cancelled, shutdown is raised")
		if e := server.Shutdown(context.Background()); e != nil {
			m.log.Emergency("graceful shutdown error, %v", logger.Args(e))
		}
	}()

	if err = server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			m.log.Emergency("server is shutdown with error, %v", logger.Args(err))
		} else {
			err = nil
		}
	}
	return
}

// New
func New(ctx context.Context, mux *chi.Mux, log logger.Logger, cfg Config, managers Managers) *Http {
	return &Http{
		ctx:     ctx,
		cfg:     cfg,
		oauth:   managers.oauth,
		session: managers.session,
		mux:     mux,
		log:     log.WithFields(logger.Fields{"service": prefix}),
	}
}
