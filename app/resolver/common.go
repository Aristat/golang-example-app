package resolver

import (
	"context"

	graphql1 "github.com/aristat/golang-example-app/generated/graphql"
	"github.com/google/uuid"
)

type commonQueryResolver struct{ *Resolver }

// CommonQuery
func (r *Resolver) CommonQuery() graphql1.CommonQueryResolver {
	return &commonQueryResolver{r}
}

// QUERY

// Common
func (r *queryResolver) Common(ctx context.Context) (*graphql1.CommonQuery, error) {
	return &graphql1.CommonQuery{}, nil
}

// UUID
func (r *commonQueryResolver) UUID(ctx context.Context, obj *graphql1.CommonQuery) (string, error) {
	id, err := uuid.NewRandom()
	return id.String(), err
}
