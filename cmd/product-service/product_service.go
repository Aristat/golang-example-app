package product_service

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"

	"github.com/aristat/golang-example-app/app/config"

	"github.com/aristat/golang-example-app/app/common"
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/opentracing/opentracing-go"

	"github.com/aristat/golang-example-app/generated/resources/proto/health_checks"
	"github.com/aristat/golang-example-app/generated/resources/proto/products"
	"github.com/golang/protobuf/ptypes/empty"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// Config
type Config struct {
	Port           string
	HealthCheckUrl string
	NatsURL        string
	Subject        string
}

type server struct {
	logger logger.Logger
	cfg    Config
	products.UnimplementedProductsServer
}

func (s *server) ListProduct(ctx context.Context, in *products.ListProductIn) (*products.ListProductOut, error) {
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

	conn, err := grpc.Dial(s.cfg.HealthCheckUrl, opts...)
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

	// test result
	out := &products.ListProductOut{Status: products.ListProductOut_OK, Products: []*products.Product{}}
	out.Products = append(out.Products, &products.Product{Id: 1, Name: "first_product"})
	out.Products = append(out.Products, &products.Product{Id: 2, Name: "second_product"})

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
			conf, c, e := config.Build()
			if e != nil {
				panic(e)
			}
			defer c()

			clientConfig := Config{}
			e = conf.UnmarshalKey("services.productService", &clientConfig)
			if e != nil {
				log.Fatal("Config initialize error")
			}

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

			tracer := common.GenerateTracerForTestClient("golang-example-app-product-service", conf)
			opentracing.SetGlobalTracer(tracer)

			log.Info("Start product service %s", logger.Args(clientConfig.Port))
			lis, err := net.Listen("tcp", ":"+clientConfig.Port)
			if err != nil {
				panic(err)
			}

			nc, err := nats.Connect(clientConfig.NatsURL)
			if err != nil {
				log.Error(err.Error())
				panic(err)
			}
			defer nc.Close()

			sc, _ := stan.Connect("test-cluster", "example-subscriber", stan.NatsConn(nc))
			natsService := natsService{logger: log}
			_, err = sc.QueueSubscribe(clientConfig.Subject, "worker", natsService.workerHanlder, stan.DurableName("i-will-remember"), stan.MaxInflight(1), stan.SetManualAckMode())

			if err != nil {
				sc.Close()
				log.Error(err.Error())
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
			products.RegisterProductsServer(s, &server{logger: log, cfg: clientConfig})

			if err := s.Serve(lis); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
}
