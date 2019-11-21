package resolver

import (
	"context"

	graphql1 "github.com/aristat/golang-example-app/generated/graphql"
	"github.com/spf13/cast"
)

type usersQueryResolver struct{ *Resolver }
type usersMutationResolver struct{ *Resolver }

func (r *Resolver) UsersMutation() graphql1.UsersMutationResolver {
	return &usersMutationResolver{r}
}
func (r *Resolver) UsersQuery() graphql1.UsersQueryResolver {
	return &usersQueryResolver{r}
}

// QUERY

func (r *queryResolver) Users(ctx context.Context) (*graphql1.UsersQuery, error) {
	return &graphql1.UsersQuery{}, nil
}

func (r *usersQueryResolver) One(ctx context.Context, obj *graphql1.UsersQuery, email string) (*graphql1.UsersOneOut, error) {
	user, err := r.repo.Users.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	userData := &graphql1.UsersOneOut{
		ID:    cast.ToString(&user.ID),
		Email: user.Email,
	}

	return userData, nil
}

// MUTATIONS

func (r *mutationResolver) Users(ctx context.Context) (*graphql1.UsersMutation, error) {
	return &graphql1.UsersMutation{}, nil
}

func (r *usersMutationResolver) CreateUser(ctx context.Context, obj *graphql1.UsersMutation, email string, password string) (*graphql1.UsersCreateOut, error) {
	user, err := r.repo.Users.CreateUser(email, password)

	if err != nil {
		return nil, err
	}

	userData := &graphql1.UsersCreateOut{
		ID:     cast.ToString(user.ID),
		Email:  user.Email,
		Status: graphql1.UsersCreateOutStatusOk,
	}

	return userData, nil
}
