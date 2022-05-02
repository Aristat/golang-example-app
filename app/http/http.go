package http

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"

	"github.com/go-session/session"

	"github.com/aristat/golang-example-app/app/logger"
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
	session *session.Manager
	log     logger.Logger
	mux     *chi.Mux
}

// ListenAndServe
func (m *Http) ListenAndServe(wg *sync.WaitGroup, bind ...string) (server *http.Server) {
	bindAddress := m.cfg.Bind

	if len(bind) > 0 && len(bind[0]) > 0 {
		bindAddress = bind[0]
	}

	server = &http.Server{
		Addr:    bindAddress,
		Handler: m.mux,
	}

	go func() {
		defer wg.Done()

		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				m.log.Emergency("Server is shutdown with error, %v", logger.Args(err))
			} else {
				err = nil
			}
		}

		m.log.Info("HTTP Server stopped successfully")
	}()

	return server
}

// New
func New(ctx context.Context, mux *chi.Mux, log logger.Logger, cfg Config) *Http {
	return &Http{
		ctx: ctx,
		cfg: cfg,
		mux: mux,
		log: log.WithFields(logger.Fields{"service": prefix}),
	}
}
