package graphql_resolver

import (
	"context"
	"fmt"

	"github.com/aristat/golang-example-app/generated/resources/proto/products"
)

type ProductServerMock struct{}

func (s *ProductServerMock) ListProduct(ctx context.Context, in *products.ListProductIn) (*products.ListProductOut, error) {
	productIds := []int64{1, 2, 3, 4, 5}
	out := &products.ListProductOut{Status: products.ListProductOut_OK, Products: []*products.Product{}}

	for _, id := range productIds {
		out.Products = append(out.Products, &products.Product{Id: id, Name: fmt.Sprintf("product_%d", id)})
	}

	return out, nil
}
