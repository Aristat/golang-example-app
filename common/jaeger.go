package common

import (
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func GenerateTracerForTestClient(serviceName string) opentracing.Tracer {
	jaegerCfg := jaegerConfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: false,
		},
	}

	tracer, _, e := jaegerCfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if e != nil {
		log.Fatal("Jaeger initialize error")
	}

	return tracer
}
