package natsmq

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"log"
	"time"
)

type NatsConn struct {
	Conn   *nats.Conn
	Tracer opentracing.Tracer
	Logger *zap.SugaredLogger
}

type NatsEncodedConn struct {
	Conn   *nats.EncodedConn
	Tracer opentracing.Tracer
	Logger *zap.SugaredLogger
}

func setupDefaultNatsConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := 10 * time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, e error) {
		log.Printf("[NATS Disconnect handler] Disconnected due: %s. Will attempt reconnects for %.0fm", e, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("[NATS reconnect handler] Successfullly reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Printf("[NATS closed handler] NATS connection closed: %v", nc.LastError())
	}))
	opts = append(opts, nats.ErrorHandler(func(nc *nats.Conn, _ *nats.Subscription, err error) {
		log.Printf("[NATS error handler] Connection error: %s", err)
	}))
	return opts
}

func NewConn(ctx context.Context, brokersAddr string, opts ...nats.Option) (*nats.Conn, error) {
	opts = setupDefaultNatsConnOptions(opts)

	nc, err := nats.Connect(brokersAddr, opts...)
	if err != nil {
		return nil, err
	}

	return nc, nil
}

func NewEncodedConn(ctx context.Context, brokersAddr string, opts ...nats.Option) (*nats.EncodedConn, error) {

	nc, err := NewConn(ctx, brokersAddr, opts...)
	if err != nil {
		return nil, err
	}

	encoded, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}
