package graphql

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/generated/graphql"
	"github.com/go-chi/chi"
	"github.com/vektah/gqlparser/gqlerror"
)

const (
	prefix = "app.graphql"
)

var errInternalServer = errors.New("internal server error")

// GraphQL
type GraphQL struct {
	ctx      context.Context
	resolver *graphql.Config
	log      logger.Logger
	cfg      Config
}

// Use
func (g *GraphQL) Use(router *chi.Mux) {
	router.Use(g.cfg.Middleware...)
}

// Routers
func (g *GraphQL) Routers(router *chi.Mux) {
	options := []handler.Option{
		handler.IntrospectionEnabled(g.cfg.Introspection),
		handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
			g.log.Alert("unhandled panic, err: %v", logger.Args(err))
			return nil
		}),
		handler.ErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
			if e != nil {
				g.log.Alert("recover on middleware, err: %v", logger.Args(e))
				goto done
			}
			e = errInternalServer
		done:
			return gqlgen.DefaultErrorPresenter(ctx, e)
		}),
	}

	if g.cfg.Debug {
		router.Handle("/", handler.Playground(g.cfg.Name, "/query"))
		options = append(options, handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
			startTime := time.Now()
			rc := gqlgen.GetRequestContext(ctx)
			resp := next(ctx)
			e := strings.ReplaceAll(rc.Errors.Error(), "\n", " ")
			g.log.Debug("\nVARS:\n%+v\nQUERY:\n%v\nRESPONSE:\n%v\nERROR:\n%v\n",
				logger.Args(rc.Variables, strings.TrimRight(rc.RawQuery, "\n"), string(resp), e),
				logger.WithFields(logger.Fields{
					"time": time.Since(startTime).String(),
				}),
			)
			return resp
		}))
	}

	router.Handle("/query",
		handler.GraphQL(
			graphql.NewExecutableSchema(*g.resolver),
			options...,
		),
	)
}

// Config
type Config struct {
	Debug         bool
	Introspection bool
	Name          string
	Middleware    []func(http.Handler) http.Handler
}

// New
func New(ctx context.Context, resolver graphql.Config, log logger.Logger, cfg Config) *GraphQL {
	log = log.WithFields(logger.Fields{"service": prefix})
	return &GraphQL{
		ctx:      ctx,
		resolver: &resolver,
		cfg:      cfg,
		log:      log,
	}
}
