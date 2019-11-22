package resolver

import (
	"context"
	"time"

	"github.com/aristat/golang-example-app/app/dataloader"

	"github.com/aristat/golang-example-app/app/db/domain"

	"github.com/aristat/golang-example-app/app/grpc"
	"github.com/aristat/golang-example-app/common"
	"github.com/aristat/golang-example-app/generated/graphql"
	"github.com/aristat/golang-example-app/generated/resources/proto/products"
	"github.com/spf13/cast"
)

type productsQueryResolver struct{ *Resolver }
type productResolver struct{ *Resolver }

func (r *Resolver) ProductsQuery() graphql.ProductsQueryResolver {
	return &productsQueryResolver{r}
}

// QUERY

func (r *queryResolver) Products(ctx context.Context) (*graphql.ProductsQuery, error) {
	return &graphql.ProductsQuery{}, nil
}

func (r *productsQueryResolver) List(ctx context.Context, obj *graphql.ProductsQuery) (*graphql.ProductsListOut, error) {
	conn, d, err := grpc.GetConnGRPC(r.pollManager, common.SrvProducts)
	defer d()

	if err != nil {
		return nil, err
	}

	c := products.NewProductsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(r.cfg.ProductTimeout))
	defer cancel()

	productOut, err := c.ListProduct(ctx, &products.ListProductIn{Id: 1})
	if err != nil {
		return nil, err
	}

	list := make([]*domain.Product, len(productOut.Products))
	for i, product := range productOut.Products {
		list[i] = &domain.Product{
			ID:   cast.ToInt(&product.Id),
			Name: product.Name,
		}
	}

	return &graphql.ProductsListOut{Products: list}, nil
}

func (r *queryResolver) ProductsRoot(ctx context.Context) ([]*domain.Product, error) {
	return []*domain.Product{{ID: 1, Name: "test1"}, {ID: 2, Name: "test2"}}, nil
}

func (r *Resolver) Product() graphql.ProductResolver {
	return &productResolver{r}
}

func (r *productResolver) ProductItems(ctx context.Context, obj *domain.Product) ([]*graphql.ProductItem, error) {
	r.log.Info("ProductItems Start Request")
	return dataloader.CtxLoaders(ctx).ProductItemsByProduct.Load(obj.ID)
}
