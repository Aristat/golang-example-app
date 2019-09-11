package product_service

import (
	"context"
	"log"
	"net"

	"github.com/uber/jaeger-client-go"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"

	"github.com/aristat/golang-example-app/generated/resources/proto/products"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) ListProduct(ctx context.Context, in *products.ListProductIn) (*products.ListProductOut, error) {
	log.Printf("Received ListProduct: %v", in.Id)

	out := &products.ListProductOut{Status: products.ListProductOut_OK, Products: []*products.Product{}}
	out.Products = append(out.Products, &products.Product{Id: 2, Name: "first_product"})

	return out, nil
}

var (
	//bind string
	Cmd = &cobra.Command{
		Use:           "product-service",
		Short:         "Product service",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			log.SetFlags(log.Lshortfile | log.LstdFlags)

			jaegerCfg := jaegerConfig.Configuration{
				ServiceName: "golang-example-app-product-service",
				Sampler: &jaegerConfig.SamplerConfig{
					Type:  "const",
					Param: 1,
				},
				Reporter: &jaegerConfig.ReporterConfig{
					LogSpans: false,
				},
			}

			tracer, _, e := jaegerCfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
			if e != nil {
				log.Fatal("Jaeger initialize error")
			}

			lis, err := net.Listen("tcp", port)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer(
				grpc.StreamInterceptor(
					grpc_middleware.ChainStreamServer(
						otgrpc.OpenTracingStreamServerInterceptor(tracer),
					),
				),
				grpc.UnaryInterceptor(
					grpc_middleware.ChainUnaryServer(
						otgrpc.OpenTracingServerInterceptor(tracer),
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
