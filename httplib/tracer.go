package httplib

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
	"log"
)

func InitTracing(ctx context.Context, host, name string) (context.Context, opentracing.Tracer, io.Closer, error) {
	span := opentracing.StartSpan(fmt.Sprintf("%s-startup", name))
	ctx = opentracing.ContextWithSpan(ctx, span)

	metrics := prometheus.New()

	traceTransport, err := jaeger.NewUDPTransport(host, 0)
	if err != nil {
		log.Printf("Unable to init jaeger udp transport: %s", err)
		return ctx, nil, nil, err
	}
	tracer, closer, err := config.Configuration{
		ServiceName: name,
	}.NewTracer(
		config.Sampler(jaeger.NewConstSampler(true)),
		config.Reporter(jaeger.NewRemoteReporter(traceTransport, jaeger.ReporterOptions.Logger(jaeger.StdLogger))),
		config.Metrics(metrics),
	)
	if err != nil {
		log.Printf("Unable to create tracer conn: %s", err)
		return ctx, nil, nil, err
	}

	return ctx, tracer, closer, nil
}
