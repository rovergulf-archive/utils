package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
	"io"
)

type Jaeger struct {
	opentracing.Tracer
	ServiceName string
	Metrics     *prometheus.Factory
	Closer      io.Closer
	logger      *zap.SugaredLogger
}

func NewJaeger(ctx context.Context, logger *zap.SugaredLogger, serviceName, address string) (*Jaeger, error) {
	j := new(Jaeger)

	j.logger = logger
	j.ServiceName = serviceName
	j.Metrics = prometheus.New()

	span := opentracing.StartSpan(fmt.Sprintf("%s startup", j.ServiceName))
	ctx = opentracing.ContextWithSpan(context.Background(), span)
	defer span.Finish()

	traceTransport, err := jaeger.NewUDPTransport(address, 0)
	if err != nil {
		j.logger.Errorf("Unable to init tracing transport: %s", err)
		return nil, err
	}

	tracer, closer, err := config.Configuration{
		ServiceName: j.ServiceName,
	}.NewTracer(
		config.Sampler(jaeger.NewConstSampler(true)),
		config.Reporter(jaeger.NewRemoteReporter(traceTransport, jaeger.ReporterOptions.Logger(jaeger.StdLogger))),
		config.Metrics(j.Metrics),
	)
	if err != nil {
		j.logger.Errorf("Unable to start tracer: %s", err)
		return nil, err
	}

	j.Tracer = tracer
	j.Closer = closer

	j.logger.Debugw("Jaeger tracing client initialized", "collector_url", address)
	return j, nil
}

func (j *Jaeger) GracefulShutdown(ctx context.Context) {
	if j.Closer != nil {
		if err := j.Closer.Close(); err != nil {
			j.logger.Errorf("Unable to close tracing io.Closer: %s", err)
		}
	}
}
