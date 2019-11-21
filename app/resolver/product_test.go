package resolver_test

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"testing"

	"github.com/aristat/golang-example-app/common"

	"github.com/aristat/golang-example-app/generated/resources/proto/products"

	"google.golang.org/grpc"

	"github.com/stretchr/testify/assert"

	grpc1 "github.com/aristat/golang-example-app/app/grpc"
	"github.com/aristat/golang-example-app/app/resolver"
	graphql1 "github.com/aristat/golang-example-app/generated/graphql"
)

var (
	productServerHost = "localhost"
	productServerPort = ":50051"
)

func TestMain(m *testing.M) {
	lis, err := net.Listen("tcp", productServerPort)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer lis.Close()

	s := grpc.NewServer()
	products.RegisterProductsServer(s, &resolver.ProductServerMock{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	code := m.Run()

	os.Exit(code)
}

func TestList(t *testing.T) {
	var opts []grpc.DialOption
	ctx := context.Background()
	opts = append(opts, grpc.WithInsecure())

	pool, _ := grpc1.NewPool(ctx, common.SrvProducts, productServerHost+productServerPort, grpc1.ConnOptions(opts...))
	grpc1.SetPool(pool, common.SrvProducts)

	cfg, _, err := resolver.BuildTest()
	if err != nil {
		assert.Failf(t, "resolver instance failed, err: %v", err.Error())
		return
	}

	obj := graphql1.ProductsQuery{}

	out, err := cfg.Resolvers.ProductsQuery().List(ctx, &obj)
	if err != nil {
		assert.Failf(t, "request failed, err: %v", err.Error())
		return
	}

	jsonProducts, _ := json.Marshal(out.Products)
	t.Log(string(jsonProducts))
	assert.Equal(t, len(out.Products), 5)
}
