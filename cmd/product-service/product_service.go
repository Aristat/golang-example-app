package product_service

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"go.opentelemetry.io/otel/propagation"

	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"

	"github.com/aristat/golang-example-app/app/config"

	"github.com/aristat/golang-example-app/app/common"
	"github.com/aristat/golang-example-app/app/logger"

	"github.com/aristat/golang-example-app/generated/resources/proto/health_checks"
	"github.com/aristat/golang-example-app/generated/resources/proto/products"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

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
	tracer := otel.GetTracerProvider()

	conn, err := grpc.Dial(s.cfg.HealthCheckUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logger.UnaryClientInterceptor(s.logger, true),
			otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(tracer), otelgrpc.WithPropagators(propagation.TraceContext{})),
		),
		grpc.WithChainStreamInterceptor(
			logger.StreamClientInterceptor(s.logger, true),
			otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(tracer), otelgrpc.WithPropagators(propagation.TraceContext{})),
		),
	)

	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := health_checks.NewHealthChecksClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	isAliveOut, err := c.IsAlive(ctx, &emptypb.Empty{})
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

			tracer, e := common.GenerateTracerForTestClient("golang-example-app-product-service", conf)
			otel.SetTracerProvider(tracer)

			if e != nil {
				panic(e)
			}

			defer func() {
				if err := tracer.Shutdown(context.Background()); err != nil {
					log.Printf("Error shutting down tracer provider: %v", err)
				}
			}()

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
				grpc.ChainUnaryInterceptor(
					logger.UnaryServerInterceptor(log, true),
					otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tracer), otelgrpc.WithPropagators(propagation.TraceContext{})),
				),
				grpc.ChainStreamInterceptor(
					logger.StreamServerInterceptor(log, true),
					otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tracer), otelgrpc.WithPropagators(propagation.TraceContext{})),
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
