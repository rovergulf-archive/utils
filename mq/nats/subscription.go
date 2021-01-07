package natsmq

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"github.com/rovergulf/utils/clog"
	"log"
)

type NatsSub struct {
	Ctx    context.Context
	Tracer opentracing.Tracer

	messages chan *nats.Msg
	errors   chan error
	quit     chan struct{}
	subject  string
	response string
	conn     *nats.EncodedConn
	sub      *nats.Subscription
}

func NewSubscriptionWithTracing(ctx context.Context, tracer opentracing.Tracer, brokerAddr, subject string, opts ...nats.Option) (*NatsSub, error) {
	if tracer == nil {
		return nil, fmt.Errorf("'tracer opentracting.Tracer' provided as nil")
	}

	sub, err := NewSubscription(ctx, brokerAddr, subject, opts...)
	if err != nil {
		return nil, err
	}

	sub.Tracer = tracer

	return sub, nil
}

func NewSubscription(ctx context.Context, brokerAddr, subject string, opts ...nats.Option) (*NatsSub, error) {
	c := new(NatsSub)

	c.Ctx = ctx
	c.messages = make(chan *nats.Msg)
	c.errors = make(chan error)
	c.quit = make(chan struct{})
	c.subject = subject

	opts = append(opts, nats.Name("chan-"+subject))

	enc, err := NewEncodedConn(ctx, brokerAddr, opts...)
	if err != nil {
		clog.Errorf("Unable to create NATS encoded connection: %s", err)
		return nil, err
	}

	c.conn = enc

	sub, err := c.conn.Subscribe(subject, func(msg *nats.Msg) {
		c.messages <- msg
	})
	if err != nil {
		clog.Errorf("Unable to start scheduler subject subscription: %s", err)
		return nil, err
	}

	c.sub = sub

	log.Printf("Initialized NATS subscription at '%s' subject\n", c.subject)
	return c, nil
}

func (ns *NatsSub) StartConsumption(ctx context.Context, handler func(data []byte) error) {
loop:
	for {
		select {
		case <-ns.Ctx.Done():
			log.Println("Received shutdown signal, stopping subscription")
			if err := ns.Stop(); err != nil { // ???
				log.Printf("Failed to unsubscribe at %s: %s", ns.subject, err)
			}
			if ns.conn != nil {
				ns.conn.Close()
			}
			break loop
		case m := <-ns.Messages():
			var span opentracing.Span

			// check if we have available handler
			delivered, _ := ns.sub.Delivered()
			log.Printf("Successfully received '%s' message with increment: %d", ns.subject, delivered)

			if ns.Tracer != nil {
				span = ns.Tracer.StartSpan(fmt.Sprintf("[%s:%d]", ns.subject, delivered))
			}

			if err := handler(m.Data); err != nil {
				log.Printf("Unable to handle nats '%s' subscription message: %s", ns.subject, err)
			}
			if len(m.Reply) > 0 {
				if err := m.Respond([]byte(m.Reply)); err != nil {
					log.Printf("Unable to respond nats message: %s", err)
				} else {
					log.Printf("Succesfully responed [%s: %s]", ns.subject, m.Reply)
				}
			}

			if span != nil {
				span.Finish()
			}
		case e := <-ns.Errors():
			log.Printf("Subscription [%s] error: %s", ns.subject, e)
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
		ns.conn.Drain()
		ns.conn.Close()
	}
	return ns.sub.Unsubscribe()
}
