package resolver

import (
	"context"

	"github.com/aristat/golang-example-app/app/grpc"

	"github.com/aristat/golang-example-app/app/db/repo"

	"github.com/aristat/golang-example-app/app/logger"
	graphql1 "github.com/aristat/golang-example-app/generated/graphql"
)

var prefix = "app.resolver"

// Config
type Config struct {
}

// Managers
type Managers struct {
	Repo        *repo.Repo
	PollManager *grpc.PoolManager
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

// Resolver config graphql resolvers
type Resolver struct {
	ctx         context.Context
	log         logger.Logger
	cfg         Config
	repo        *repo.Repo
	pollManager *grpc.PoolManager
}

// Mutation returns root graphql mutation resolver
func (r *Resolver) Mutation() graphql1.MutationResolver {
	return &mutationResolver{r}
}

// Query returns root graphql query resolver
func (r *Resolver) Query() graphql1.QueryResolver {
	return &queryResolver{r}
}

// New
func New(ctx context.Context, log logger.Logger, cfg Config, managers Managers) graphql1.Config {
	log = log.WithFields(logger.Fields{"service": prefix})
	c := graphql1.Config{
		Resolvers: &Resolver{
			ctx:         ctx,
			log:         log,
			cfg:         cfg,
			repo:        managers.Repo,
			pollManager: managers.PollManager,
		},
	}
	return c
}
