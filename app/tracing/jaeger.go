package tracing

import (
	"context"

	"github.com/aristat/golang-example-app/app/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go/config"
)

const prefix = "app.tracer"

type (
	Tracer = opentracing.Tracer
)

// New returns instance implemented of opentracing.Tracer interface
func New(ctx context.Context, log logger.Logger, cfg config.Configuration, option ...config.Option) (Tracer, error) {
	log = log.WithFields(logger.Fields{"service": prefix})
	tracer, closer, e := cfg.NewTracer(option...)
	if e != nil {
		return tracer, errors.WithMessage(e, prefix)
	}
	go func() {
		<-ctx.Done()
		if e := closer.Close(); e != nil {
			log.Error("%v", logger.Args(e))
		}
	}()
	return tracer, errors.WithMessage(e, prefix)
}
