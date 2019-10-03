package resolver

import (
	"context"
	"time"

	"github.com/aristat/golang-example-app/app/grpc"
	"github.com/aristat/golang-example-app/common"
	"github.com/aristat/golang-example-app/generated/resources/proto/products"
	"github.com/spf13/cast"

	graphql1 "github.com/aristat/golang-example-app/generated/graphql"
)

type productsQueryResolver struct{ *Resolver }

func (r *Resolver) ProductsQuery() graphql1.ProductsQueryResolver {
	return &productsQueryResolver{r}
}

// QUERY

func (r *queryResolver) Products(ctx context.Context) (*graphql1.ProductsQuery, error) {
	return &graphql1.ProductsQuery{}, nil
}

func (r *productsQueryResolver) List(ctx context.Context, obj *graphql1.ProductsQuery) (*graphql1.ProductsListOut, error) {
	conn, d, err := grpc.GetConnGRPC(r.pollManager, common.SrvProducts)
	defer d()

	if err != nil {
		return nil, err
	}

	c := products.NewProductsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	productOut, err := c.ListProduct(ctx, &products.ListProductIn{Id: 1})
	if err != nil {
		return nil, err
	}

	list := make([]*graphql1.Product, len(productOut.Products))
	for i, product := range productOut.Products {
		list[i] = &graphql1.Product{
			ID:   cast.ToString(&product.Id),
			Name: product.Name,
		}
	}

	return &graphql1.ProductsListOut{Products: list}, nil
}
