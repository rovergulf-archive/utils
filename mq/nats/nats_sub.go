package natsmq

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"time"
)

type NatsSub struct {
	Tracer opentracing.Tracer
	Logger *zap.SugaredLogger

	messages chan *nats.Msg
	errors   chan error
	quit     chan struct{}
	subject  string
	name     string
	response string
	conn     *nats.EncodedConn
	sub      *nats.Subscription
}

func NewSubscriptionWithTracing(c *NatsSubOpts, tracer opentracing.Tracer) (*NatsSub, error) {
	if tracer == nil {
		return nil, fmt.Errorf("'tracer opentracting.Tracer' provided as nil")
	}

	sub, err := NewSubscription(c)
	if err != nil {
		return nil, err
	}
	sub.Tracer = tracer

	return sub, nil
}

func NewSubscription(c *NatsSubOpts) (*NatsSub, error) {
	ns := &NatsSub{
		Tracer:   c.Tracer,
		Logger:   c.Logger.Named(c.Subject),
		messages: make(chan *nats.Msg),
		errors:   make(chan error),
		quit:     make(chan struct{}),
		subject:  c.Subject,
		name:     fmt.Sprintf("chan-%s-%d", c.Subject, time.Now().Unix()),
		response: c.Subject + "-sub",
	}

	c.NatsConn = append(c.NatsConn, nats.Name(ns.name))

	enc, err := NewEncodedConn(c.Config)
	if err != nil {
		ns.Logger.Errorf("Unable to create NATS encoded connection: %s", err)
		return nil, err
	}

	ns.conn = enc

	sub, err := ns.conn.Subscribe(ns.subject, func(msg *nats.Msg) {
		ns.messages <- msg
	})
	if err != nil {
		ns.Logger.Errorf("Unable to start scheduler subject subscription: %s", err)
		return nil, err
	}

	ns.sub = sub

	ns.Logger.Infof("Initialized NATS subscription at '%s' subject", ns.subject)
	return ns, nil
}

type NatsSubHandler func(data []byte, reply string) error

func (ns *NatsSub) StartConsumption(ctx context.Context, handler NatsSubHandler) {
loop:
	for {
		select {
		case <-ctx.Done():
			ns.Logger.Info("Received shutdown signal, stopping subscription")
			if err := ns.Stop(); err != nil { // ???
				ns.Logger.Infof("Failed to unsubscribe at %s: %s", ns.subject, err)
			}
			if ns.conn != nil {
				ns.conn.Close()
			}
			break loop
		case m := <-ns.Messages():
			var span opentracing.Span

			if ns.Tracer != nil {
				span = ns.Tracer.StartSpan(ns.subject + "-event")
				span.SetTag("subject", ns.subject)
				span.SetTag("m_reply", m.Reply)
				ctx = opentracing.ContextWithSpan(ctx, span)
			}

			if err := handler(m.Data, m.Reply); err != nil {
				ns.Logger.Infof("Unable to handle nats '%s' subscription message: %s", ns.subject, err)
			}
			if len(m.Reply) > 0 {
				if err := m.Respond([]byte(m.Reply)); err != nil {
					ns.Logger.Infof("Unable to respond nats message: %s", err)
				} else {
					ns.Logger.Infof("Succesfully responed [%s: %s]", ns.subject, m.Reply)
				}
			}

			// check if we have available handler
			delivered, _ := ns.sub.Delivered()
			ns.Logger.Infof("Successfully received '%s' message with increment: %d", ns.subject, delivered)

			if span != nil {
				span.Finish()
			}
		case e := <-ns.Errors():
			ns.Logger.Infof("Subscription [%s] error: %s", ns.subject, e)
		}
	}
}

func (ns *NatsSub) Messages() <-chan *nats.Msg {
	return ns.messages
}

func (ns *NatsSub) Errors() <-chan error {
	return ns.errors
}

func (ns *NatsSub) Stop() error {
	if ns.conn != nil {
		ns.conn.Close()
	}
	return ns.sub.Unsubscribe()
}
