package common

import (
	"log"

	"github.com/spf13/viper"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func GenerateTracerForTestClient(serviceName string, cfg *viper.Viper) opentracing.Tracer {
	jaegerCfg := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: false,
		},
	}

	e := cfg.UnmarshalKey("tracing.jaeger", &jaegerCfg)
	if e != nil {
		log.Fatal("Jaeger initialize config error")
	}

	// redefine for test services
	jaegerCfg.ServiceName = serviceName

	tracer, _, e := jaegerCfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if e != nil {
		log.Fatal("Jaeger initialize client error")
	}

	return tracer
}
