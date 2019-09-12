package product_service

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/aristat/golang-example-app/common"
	"github.com/opentracing/opentracing-go"

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

type server struct {
	logger logger.Logger
}

func (s *server) ListProduct(ctx context.Context, in *products.ListProductIn) (*products.ListProductOut, error) {
	s.logger.Printf("Received ListProduct: %v", in.Id)

	tracer := opentracing.GlobalTracer()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			logger.UnaryClientInterceptor(s.logger, true),
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
		)))
	opts = append(opts, grpc.WithStreamInterceptor(
		grpc_middleware.ChainStreamClient(
			logger.StreamClientInterceptor(s.logger, true),
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
			log, c, e := logger.Build()
			if e != nil {
				panic(e)
			}
			defer c()

			defer func() {
				if r := recover(); r != nil {
					if re, _ := r.(error); re != nil {
						log.Error(re.Error())
					} else {
						log.Alert("unhandled panic, err: %v", logger.Args(r))
					}
				}
			}()

			tracer := common.GenerateTracerForTestClient("golang-example-app-product-service")
			opentracing.SetGlobalTracer(tracer)

			lis, err := net.Listen("tcp", port)
			if err != nil {
				panic(err)
			}

			s := grpc.NewServer(
				grpc.UnaryInterceptor(
					grpc_middleware.ChainUnaryServer(
						logger.UnaryServerInterceptor(log, true),
						grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
					),
				),
				grpc.StreamInterceptor(
					grpc_middleware.ChainStreamServer(
						logger.StreamServerInterceptor(log, true),
						grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
					),
				),
			)
			products.RegisterProductsServer(s, &server{logger: log})

			if err := s.Serve(lis); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
}
