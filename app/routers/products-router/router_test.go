package products_router_test

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	products_router "github.com/aristat/golang-example-app/app/routers/products-router"
	"github.com/stretchr/testify/assert"

	"github.com/aristat/golang-example-app/app/common"
	grpc1 "github.com/aristat/golang-example-app/app/grpc"

	"github.com/aristat/golang-example-app/app/resolver"

	"github.com/aristat/golang-example-app/generated/resources/proto/products"

	"google.golang.org/grpc"
)

var (
	productServerHost = "localhost"
	grpcPort          string
)

func TestMain(m *testing.M) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer lis.Close()

	grpcPort = ":" + strconv.Itoa(lis.Addr().(*net.TCPAddr).Port)
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

func TestGetProductsGrpc(t *testing.T) {
	var opts []grpc.DialOption
	ctx := context.Background()
	opts = append(opts, grpc.WithInsecure())

	pool, _ := grpc1.NewPool(ctx, common.SrvProducts, productServerHost+grpcPort, grpc1.ConnOptions(opts...))
	grpc1.SetPool(pool, common.SrvProducts)

	tests := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "successful",
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider, _, e := products_router.BuildTest()
			assert.Nil(t, e, "err should be nil")

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/products_grpc", strings.NewReader(""))
			provider.Router.GetProductsGrpc(rec, req)

			assert.Equal(t, test.expectedCode, rec.Code)
			assert.NotNil(t, rec.Body)
		})
	}
}
