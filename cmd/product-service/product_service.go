package product_service

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/aristat/golang-example-app/common"

	"github.com/aristat/golang-example-app/generated/resources/proto/health_checks"
	"github.com/aristat/golang-example-app/generated/resources/proto/products"
	"github.com/golang/protobuf/ptypes/empty"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) ListProduct(ctx context.Context, in *products.ListProductIn) (*products.ListProductOut, error) {
	log.Printf("Received ListProduct: %v", in.Id)

	tracer := opentracing.GlobalTracer()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)))
	opts = append(opts, grpc.WithStreamInterceptor(
		grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)))

	conn, err := grpc.Dial("localhost:50052", opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := health_checks.NewHealthChecksClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	isAliveOut, err := c.IsAlive(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	if isAliveOut.Status != health_checks.IsAliveOut_OK {
		return nil, errors.New("Heal checks not working")
	}

	out := &products.ListProductOut{Status: products.ListProductOut_OK, Products: []*products.Product{}}
	out.Products = append(out.Products, &products.Product{Id: 2, Name: "first_product"})

	return out, nil
}

// Example service, which gives some data
var (
	//bind string
	Cmd = &cobra.Command{
		Use:           "product-service",
		Short:         "Product service",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			log.SetFlags(log.Lshortfile | log.LstdFlags)

			tracer := common.GenerateTracerForTestClient("golang-example-app-product-service")
			opentracing.SetGlobalTracer(tracer)

			lis, err := net.Listen("tcp", port)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer(
				grpc.StreamInterceptor(
					grpc_middleware.ChainStreamServer(
						grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
					),
				),
				grpc.UnaryInterceptor(
					grpc_middleware.ChainUnaryServer(
						grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
					),
				),
			)
			products.RegisterProductsServer(s, &server{})
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		},
	}
)

func init() {
}
