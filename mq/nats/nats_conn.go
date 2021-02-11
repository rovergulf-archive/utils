package natsmq

import (
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
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

func (c *NatsConn) setupDefaultNatsConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := 10 * time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, e error) {
		c.Logger.Warnf("[NATS Disconnect handler] Disconnected due: %s. Will attempt reconnects for %.0fm", e, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		c.Logger.Warnf("[NATS reconnect handler] Successfullly reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		c.Logger.Warnf("[NATS closed handler] NATS connection closed: %v", nc.LastError())
	}))
	opts = append(opts, nats.ErrorHandler(func(nc *nats.Conn, _ *nats.Subscription, err error) {
		c.Logger.Errorf("[NATS error handler] Connection error: %s", err)
	}))

	return opts
}

func setupDefaultNatsConnOptions(lg *zap.SugaredLogger, opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := 10 * time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, e error) {
		lg.Warnf("[NATS Disconnect handler] Disconnected due: %s. Will attempt reconnects for %.0fm", e, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		lg.Warnf("[NATS reconnect handler] Successfullly reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		lg.Warnf("[NATS closed handler] NATS connection closed: %v", nc.LastError())
	}))
	opts = append(opts, nats.ErrorHandler(func(nc *nats.Conn, _ *nats.Subscription, err error) {
		lg.Errorf("[NATS error handler] Connection error: %s", err)
	}))

	return opts
}

func NewConn(lg *zap.SugaredLogger, brokersAddr string, opts ...nats.Option) (*nats.Conn, error) {
	opts = setupDefaultNatsConnOptions(lg, opts)

	nc, err := nats.Connect(brokersAddr, opts...)
	if err != nil {
		return nil, err
	}

	return nc, nil
}

func NewEncodedConn(lg *zap.SugaredLogger, brokersAddr string, opts ...nats.Option) (*nats.EncodedConn, error) {

	nc, err := NewConn(lg, brokersAddr, opts...)
	if err != nil {
		return nil, err
	}

	encoded, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}
