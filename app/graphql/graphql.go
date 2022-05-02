package graphql

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/generated/graphql"
	"github.com/go-chi/chi/v5"
	"github.com/vektah/gqlparser/v2/gqlerror"
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
func (g *GraphQL) Routers(router chi.Router) {
	srv := handler.New(graphql.NewExecutableSchema(*g.resolver))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	if g.cfg.Introspection {
		srv.Use(extension.Introspection{})
	}
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		g.log.Alert("unhandled panic, err: %v", logger.Args(err))
		return nil
	})
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		if e != nil {
			g.log.Alert("recover on middleware, err: %v", logger.Args(e))
			goto done
		}
		e = errInternalServer
	done:
		return gqlgen.DefaultErrorPresenter(ctx, e)
	})

	if g.cfg.Debug {
		srv.AroundResponses(func(ctx context.Context, next gqlgen.ResponseHandler) *gqlgen.Response {
			startTime := time.Now()
			rc := gqlgen.GetOperationContext(ctx)
			resp := next(ctx)
			g.log.Debug("\nVARS:\n%+v\nQUERY:\n%v\nRESPONSE:\n%v\nERROR:\n%v\n",
				logger.Args(rc.Variables, strings.TrimRight(rc.RawQuery, "\n"), string(resp.Data), resp.Errors),
				logger.WithFields(logger.Fields{
					"time": time.Since(startTime).String(),
				}),
			)
			return resp
		})
	}

	router.Handle("/query", srv)
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
