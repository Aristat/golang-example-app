package health_check_service

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel"

	"github.com/aristat/golang-example-app/app/config"

	"github.com/aristat/golang-example-app/app/logger"

	"github.com/aristat/golang-example-app/app/common"

	"github.com/aristat/golang-example-app/generated/resources/proto/health_checks"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// Config
type Config struct {
	Port          string
	RandomDisable bool
}

type server struct {
	cfg Config
	health_checks.UnimplementedHealthChecksServer
}

func (s *server) IsAlive(ctx context.Context, empty *emptypb.Empty) (*health_checks.IsAliveOut, error) {
	if s.cfg.RandomDisable {
		return &health_checks.IsAliveOut{Status: health_checks.IsAliveOut_OK}, nil
	}

	var status health_checks.IsAliveOut_Status

	rand.Seed(time.Now().UTC().UnixNano())
	number := rand.Intn(2-0) + 0

	if number == 1 {
		status = health_checks.IsAliveOut_OK
	} else {
		status = health_checks.IsAliveOut_NOT_OK
	}

	return &health_checks.IsAliveOut{Status: status}, nil
}

// Example service, which need for testing jaeger and grpc pool
var (
	//bind string
	Cmd = &cobra.Command{
		Use:           "health-check",
		Short:         "Health check",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			conf, c, e := config.Build()
			if e != nil {
				panic(e)
			}
			defer c()

			clientConfig := Config{}
			e = conf.UnmarshalKey("services.healthCheckService", &clientConfig)
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

			tracer, e := common.GenerateTracerForTestClient("golang-example-app-health-check-service", conf)
			otel.SetTracerProvider(tracer)

			if e != nil {
				panic(e)
			}

			defer func() {
				if err := tracer.Shutdown(context.Background()); err != nil {
					log.Printf("Error shutting down tracer provider: %v", err)
				}
			}()

			lis, err := net.Listen("tcp", ":"+clientConfig.Port)
			if err != nil {
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
			health_checks.RegisterHealthChecksServer(s, &server{cfg: clientConfig})
			if err := s.Serve(lis); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
}
