package tracing

import (
	"context"

	jaegerConfig "github.com/uber/jaeger-client-go/config"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/google/wire"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
)

// ProviderCfg
func ProviderCfg(cfg *viper.Viper) (jaegerConfig.Configuration, func(), error) {
	c := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: false,
		},
	}
	e := cfg.UnmarshalKey("tracing.jaeger", &c)
	return c, func() {}, e
}

// Provider
func Provider(ctx context.Context, cfg jaegerConfig.Configuration, log logger.Logger) (Tracer, func(), error) {
	t, e := newJaegerTracer(ctx, log, cfg, jaegerConfig.Logger(jaeger.StdLogger))
	opentracing.SetGlobalTracer(t)
	return t, func() {}, e
}

// ProviderTest
func ProviderTest() (Tracer, func(), error) {
	m := mocktracer.New()
	return m, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, ProviderCfg)
	ProviderTestSet       = wire.NewSet(ProviderTest)
)
