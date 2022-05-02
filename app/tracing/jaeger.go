package tracing

import (
	"context"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/pkg/errors"
)

type Configuration struct {
	AgentHost   string
	AgentPort   string
	ServiceName string
}

const prefix = "app.tracer"

func newJaegerTracer(ctx context.Context, configuration *Configuration, log logger.Logger) (*tracesdk.TracerProvider, error) {
	log = log.WithFields(logger.Fields{"service": prefix})

	exp, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(configuration.AgentHost),
			jaeger.WithAgentPort(configuration.AgentPort),
		),
	)
	if err != nil {
		return nil, errors.WithMessage(err, prefix)
	}

	tracer := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(configuration.ServiceName),
		)),
	)

	go func() {
		<-ctx.Done()
		if err := tracer.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	return tracer, errors.WithMessage(err, prefix)
}
