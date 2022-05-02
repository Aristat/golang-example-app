package tracing

import (
	"context"

	"go.opentelemetry.io/otel/propagation"

	"github.com/google/wire"
	"github.com/spf13/viper"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel"

	"github.com/aristat/golang-example-app/app/logger"
)

// Cfg
func Cfg(cfg *viper.Viper) (*Configuration, func(), error) {
	c := Configuration{}
	e := cfg.UnmarshalKey("tracing.jaeger", &c)
	if e != nil {
		return nil, func() {}, e
	}
	return &c, func() {}, nil
}

// Provider
func Provider(ctx context.Context, configuration *Configuration, log logger.Logger) (*tracesdk.TracerProvider, func(), error) {
	t, e := newJaegerTracer(ctx, configuration, log)
	otel.SetTracerProvider(t)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return t, func() {}, e
}

// ProviderTest
func ProviderTest() (*tracesdk.TracerProvider, func(), error) {
	m := tracesdk.NewTracerProvider()
	return m, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, Cfg)
	ProviderTestSet       = wire.NewSet(ProviderTest)
)
