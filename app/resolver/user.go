package resolver

import (
	"context"

	graphql1 "github.com/aristat/golang-example-app/generated/graphql"
	"github.com/spf13/cast"
)

func (r *queryResolver) User(ctx context.Context, email string) (*graphql1.User, error) {
	user, err := r.repo.Users.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	userData := &graphql1.User{
		ID:    cast.ToString(&user.ID),
		Email: &user.Email,
	}

	return userData, nil
}
