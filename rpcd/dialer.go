package rpcd

import (
	"context"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type RpcDialer interface {
	WithCancel(ctx context.Context) func()
	GracefulShutdown(ctx context.Context)
}

type ClientConn interface {
	Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error
	Close() error
}

type Dialer struct {
	ServiceName string
	Address     string
	Conn        *grpc.ClientConn
	Logger      *zap.SugaredLogger
}

func NewDialer(ctx context.Context, lg *zap.SugaredLogger, name, addr string, opts ...DialOption) (*Dialer, error) {
	d := new(Dialer)
	d.ServiceName = name
	d.Address = addr
	d.Logger = lg

	conn, err := Dial(addr, opts...)
	if err != nil {
		d.Logger.Errorf("Unable to dial %s: %s", addr, err)
		return nil, err
	}
	d.Conn = conn

	return d, nil
}

// DialOption allows optional config for dialer
type DialOption func(name string) (grpc.DialOption, error)

// WithTracer traces rpc calls
func WithTracer(t opentracing.Tracer) DialOption {
	return func(name string) (grpc.DialOption, error) {
		return grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(t)), nil
	}
}

// WithCredentials
func WithInsecure() DialOption {
	return func(name string) (grpc.DialOption, error) {
		return grpc.WithInsecure(), nil
	}
}

// Dial returns a load balanced grpc client conn with tracing interceptor
func Dial(addr string, opts ...DialOption) (*grpc.ClientConn, error) {
	var dialopts []grpc.DialOption

	for _, fn := range opts {
		opt, err := fn(addr)
		if err != nil {
			return nil, err
		}
		dialopts = append(dialopts, opt)
	}

	conn, err := grpc.Dial(addr, dialopts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (d *Dialer) WithCancel(ctx context.Context) func() {
	return func() {
		d.GracefulShutdown(ctx)
	}
}

func (d *Dialer) GracefulShutdown(ctx context.Context) {
	if d.Conn != nil {
		if err := d.Conn.Close(); err != nil {
			d.Logger.Errorf("Unable to close gRPC Client Connection: %s", err)
		}
	}
}
