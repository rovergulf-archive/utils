package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/rovergulf/utils/clog"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"io"
	"log"
	"os"
)

type Jaeger struct {
	opentracing.Tracer
	ServiceName string
	Metrics     *prometheus.Factory
	Closer      io.Closer
}

func NewJaeger(ctx context.Context, serviceName string) *Jaeger {
	j := new(Jaeger)

	j.ServiceName = serviceName
	j.Metrics = prometheus.New()

	return j
}

func InitJaeger(ctx context.Context, service, host string) (context.Context, *Jaeger, error) {

	ctx, tracer, closer, err := RunJaegerTracing(ctx, service, host)
	if err != nil {
		return ctx, nil, err
	}

	j := new(Jaeger)

	j.Tracer = tracer
	j.Closer = closer

	return ctx, j, nil
}

func (j *Jaeger) Start(ctx context.Context, host string) (context.Context, error) {
	ctx, tracer, closer, err := RunJaegerTracing(ctx, j.ServiceName, host)
	if err != nil {
		return ctx, err
	}

	j.Tracer = tracer
	j.Closer = closer

	log.Println("Successfully started Jaeger tracing")
	return ctx, nil
}

func (j *Jaeger) GracefulShutdown(ctx context.Context) {
	if j.Closer != nil {
		if err := j.Closer.Close(); err != nil {
			clog.Errorf("Unable to close tracing io.Closer: %s", err)
		}
	}
}

func RunJaegerTracing(ctx context.Context, serviceName, host string) (context.Context, opentracing.Tracer, io.Closer, error) {
	span := opentracing.StartSpan(fmt.Sprintf("%s startup", serviceName))
	ctx = opentracing.ContextWithSpan(context.Background(), span)
	defer span.Finish()

	envEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if len(envEndpoint) > 0 && envEndpoint != " " {
		host = envEndpoint
	}

	envHost := os.Getenv("JAEGER_AGENT_HOST")
	envPort := os.Getenv("JAEGER_AGENT_PORT")
	if len(envHost) > 0 && len(envPort) > 0 {
		host = fmt.Sprintf("%s:%s", envHost, envPort)
	}

	metrics := prometheus.New()

	traceTransport, err := jaeger.NewUDPTransport(host, 0)
	if err != nil {
		clog.Errorf("Unable to init tracing transport: %s", err)
	}
	tracer, closer, err := config.Configuration{
		ServiceName: serviceName,
	}.NewTracer(
		config.Sampler(jaeger.NewConstSampler(true)),
		config.Reporter(jaeger.NewRemoteReporter(traceTransport, jaeger.ReporterOptions.Logger(jaeger.StdLogger))),
		config.Metrics(metrics),
	)
	if err != nil {
		clog.Error(err)
		return ctx, nil, nil, err
	}

	log.Printf("[%s] Tracing enabled", serviceName)
	return ctx, tracer, closer, nil
}
