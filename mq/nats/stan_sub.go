package natsmq

import (
	"context"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"strings"
	"time"
)

type StanSub struct {
	tracer   opentracing.Tracer
	logger   *zap.SugaredLogger
	conn     *StanConn
	sub      stan.Subscription
	messages chan *stan.Msg
	errors   chan error
	quit     chan struct{}
	channel  string
}

type MessageInfo struct {
	Nuid      string `json:"nuid" yaml:"nuid"`
	Channel   string `json:"channel" yaml:"channel"`
	Sequence  uint64 `json:"sequence" yaml:"sequence"`
	Timestamp int64  `json:"timestamp" yaml:"timestamp"`
}

// NewChanSub creates connection with channel-named clientID
// creating subscription with a whole service lifetime context
func NewChanSub(c *StanSubOpts) (*StanSub, error) {
	ns := &StanSub{
		tracer:   c.Tracer,
		logger:   c.Logger.Named(c.Channel),
		messages: make(chan *stan.Msg),
		errors:   make(chan error),
		quit:     make(chan struct{}),
		channel:  c.Channel,
	}

	ch := strings.Split(ns.channel, ",")
	c.ClientId = fmt.Sprintf("%s-chan-%d", strings.Join(ch, "-"), time.Now().Unix())

	conn, err := NewStanConn(c.Config)
	if err != nil {
		ns.logger.Errorw("Unable to setup nats connection",
			"client_id", c.ClientId, "broker", c.Broker, "err", err)
		return nil, err
	}
	ns.conn = conn

	// set subscription options for fault tolerance
	// only ack manually
	c.Opts = append(c.Opts, stan.SetManualAckMode())
	c.Opts = append(c.Opts, stan.AckWait(60*time.Second))

	sub, err := ns.conn.client.Subscribe(ns.channel, func(msg *stan.Msg) {
		ns.messages <- msg
	}, c.Opts...)
	if err != nil {
		ns.logger.Errorw("Unable to subscribe",
			"client_id", c.ClientId, "chan", ns.channel, "err", err)
		return nil, err
	}
	ns.sub = sub

	ns.logger.Infow("Nats-streaming subscription started",
		"client_id", c.ClientId, "chan", ns.channel)
	return ns, nil
}

type StanSubHandler func(data []byte, info *MessageInfo) error

func (ns *StanSub) StartConsumption(ctx context.Context, handler StanSubHandler) {
loop:
	for {
		select {
		case <-ctx.Done():
			ns.logger.Infow("Received shutdown signal, stopping subscription",
				"client_id", ns.conn.clientId, "chan", ns.channel)
			ns.Stop()
			break loop
		case msg := <-ns.messages:
			var span opentracing.Span
			if ns.tracer != nil {
				span = ns.tracer.StartSpan(ns.channel + "-event")
				span.SetTag("sequence", msg.Sequence)
				span.SetTag("channel", ns.channel)
				ctx = opentracing.ContextWithSpan(ctx, span)
			}

			info := &MessageInfo{
				Nuid:      ns.conn.nuid.Next(),
				Sequence:  msg.Sequence,
				Timestamp: msg.Timestamp,
				Channel:   ns.channel,
			}

			if err := handler(msg.Data, info); err != nil {
				ns.logger.Errorw("Unable to handle nats message: %s", "chan", ns.channel, "err", err)
				ns.errors <- err
			}

			if err := msg.Ack(); err != nil {
				ns.logger.Infow("Unable to ack nats message",
					"chan", ns.channel, "seq", msg.Sequence, "err", err)
			} else {
				ns.conn.nuid.RandomizePrefix()
				ns.logger.Infow("Ack message",
					"chan", ns.channel, "seq", msg.Sequence, "nuid", info.Nuid)
			}

			if span != nil {
				span.Finish()
			}
		case e := <-ns.errors:
			ns.logger.Errorw("Subscription error", "chan", ns.channel, "err", e)
		}
	}
}

func (ns *StanSub) Messages() <-chan *stan.Msg {
	return ns.messages
}

func (ns *StanSub) Errors() <-chan error {
	return ns.errors
}

func (ns *StanSub) Stop() {
	if ns.sub != nil {
		if err := ns.sub.Unsubscribe(); err != nil {
			ns.logger.Errorw("Unable to unsubscribe",
				"chan", ns.channel, "err", err)
		}
	}

	if ns.conn != nil {
		ns.conn.Stop()
	}
}
