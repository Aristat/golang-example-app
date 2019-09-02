package resolver

import (
	"context"

	"github.com/aristat/golang-example-app/app/db/repo"

	"github.com/aristat/golang-example-app/app/logger"
	graphql1 "github.com/aristat/golang-example-app/generated/graphql"
)

var prefix = "resolver"

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

// Resolver config graphql resolvers
type Resolver struct {
	ctx  context.Context
	log  logger.Logger
	cfg  Config
	repo *repo.Repo
}

func (r *Resolver) Mutation() graphql1.MutationResolver {
	return &mutationResolver{r}
}

// Query returns root graphql query resolver
func (r *Resolver) Query() graphql1.QueryResolver {
	return &queryResolver{r}
}

func New(ctx context.Context, log logger.Logger, cfg Config, managers Managers) graphql1.Config {
	log = log.WithFields(logger.Fields{"service": prefix})
	c := graphql1.Config{
		Resolvers: &Resolver{
			ctx:  ctx,
			log:  log,
			cfg:  cfg,
			repo: managers.Repo,
		},
	}
	return c
}
