package resolver

import (
	"context"
	"time"

	"github.com/aristat/golang-example-app/generated/resources/proto/products"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	"github.com/aristat/golang-example-app/app/logger"
	graphql1 "github.com/aristat/golang-example-app/generated/graphql"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cast"
	"google.golang.org/grpc"
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
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts,
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			logger.UnaryClientInterceptor(r.log, true),
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
		)))
	opts = append(opts, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
		logger.StreamClientInterceptor(r.log, true),
		grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
	)))

	conn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	c := products.NewProductsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	productOut, err := c.ListProduct(ctx, &products.ListProductIn{Id: 1})
	if err != nil {
		return nil, err
	}
	r.log.Info("Products %v", logger.Args(productOut.Products))

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
