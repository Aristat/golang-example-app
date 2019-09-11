package health_check_service

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/aristat/golang-example-app/common"

	"github.com/aristat/golang-example-app/generated/resources/proto/health_checks"
	"github.com/golang/protobuf/ptypes/empty"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type server struct{}

func (s *server) IsAlive(ctx context.Context, empty *empty.Empty) (*health_checks.IsAliveOut, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	number := rand.Intn(2-0) + 0

	var status health_checks.IsAliveOut_Status
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
			log.SetFlags(log.Lshortfile | log.LstdFlags)

			tracer := common.GenerateTracerForTestClient("golang-example-app-health-check-service")

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
			health_checks.RegisterHealthChecksServer(s, &server{})
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		},
	}
)

func init() {
}
