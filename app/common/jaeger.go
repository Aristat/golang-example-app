package common

import (
	"github.com/aristat/golang-example-app/app/tracing"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/spf13/viper"
)

func GenerateTracerForTestClient(serviceName string, cfg *viper.Viper) (*tracesdk.TracerProvider, error) {
	configuration := tracing.Configuration{}
	err := cfg.UnmarshalKey("tracing.jaeger", &configuration)
	if err != nil {
		return nil, err
	}

	exp, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(configuration.AgentHost),
			jaeger.WithAgentPort(configuration.AgentPort),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	return tp, nil
}
