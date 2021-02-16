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
	Tracer   opentracing.Tracer
	Logger   *zap.SugaredLogger
	conn     *StanConn
	Sub      stan.Subscription
	messages chan *stan.Msg
	errors   chan error
	quit     chan struct{}
	channel  string
}

func NewChanSubWithTracer(c *StanSubOpts, t opentracing.Tracer) (*StanSub, error) {
	ns, err := NewChanSub(c)
	if err != nil {
		return nil, err
	}

	ns.Tracer = t
	ns.conn.Tracer = t
	return ns, nil
}

// NewChanSub creates connection with channel-named clientID
// creating subscription with a whole service lifetime context
func NewChanSub(c *StanSubOpts) (*StanSub, error) {
	ns := new(StanSub)
	ns.Logger = c.Logger.Named("stan-sub-" + c.Channel)
	ns.messages = make(chan *stan.Msg)
	ns.errors = make(chan error)
	ns.quit = make(chan struct{})
	ns.channel = c.Channel

	ch := strings.Split(c.Channel, ",")
	clientId := fmt.Sprintf("%s-chan-%d", strings.Join(ch, "-"), time.Now().Unix())

	conn, err := NewStanConn(c.Config)
	if err != nil {
		ns.Logger.Errorf("Unable to connect [%s:%s]: %s", clientId, c.Broker, err)
		return nil, err
	}
	ns.conn = conn

	// set subscription options for fault tolerance
	// only ack manually
	c.Opts = append(c.Opts, stan.SetManualAckMode())
	c.Opts = append(c.Opts, stan.AckWait(60*time.Second))
	//c.Opts = append(c.Opts, stan.StartWithLastReceived())

	sub, err := ns.conn.client.Subscribe(ns.channel, func(msg *stan.Msg) {
		ns.messages <- msg
	}, c.Opts...)
	if err != nil {
		ns.Logger.Errorf("Unable to subscribe [%s: %s]: %s", clientId, ns.channel, err)
		return nil, err
	}
	ns.Sub = sub

	delivered, err := sub.Delivered()
	if err != nil {
		ns.Logger.Errorf("Unable to check subscrpition delivered count: %s", err)
		return nil, err
	}

	ns.Logger.Infof("[%s] NATS subscription started for '%s' awating at delivered count %d", clientId, ns.channel, delivered)
	return ns, nil
}

type StanSubHandler func(data []byte, sequence uint64, reply string) error

func (ns *StanSub) StartConsumption(ctx context.Context, handler StanSubHandler) {
loop:
	for {
		select {
		case <-ctx.Done():
			ns.Logger.Infof("[%s] Received shutdown signal, stopping '%s' subscription", ns.conn.clientId, ns.channel)
			ns.Stop()
			break loop
		case msg := <-ns.messages:
			var span opentracing.Span
			if ns.Tracer != nil {
				span = ns.Tracer.StartSpan(ns.channel)
				span.SetTag("sequence", msg.Sequence)
				ctx = opentracing.ContextWithSpan(ctx, span)
			}

			if err := handler(msg.Data, msg.Sequence, msg.Reply); err != nil {
				ns.Logger.Errorf("Unable to handle nats '%s' subscription message: %s", ns.channel, err)
				ns.errors <- err
			} else {
				ns.Logger.Infof("[%s] Successfully received sequenced message: %d at %d", ns.channel, msg.Sequence, msg.Timestamp)
			}

			if err := msg.Ack(); err != nil {
				ns.Logger.Infof("Unable to respond nats message: %s", err)
			} else {
				ns.Logger.Infow("Successfully acked",
					"channel", ns.channel, "sequence", msg.Sequence, "reply", msg.Reply)
			}

			if span != nil {
				span.Finish()
			}
		case e := <-ns.errors:
			ns.Logger.Errorf("Subscription error: %s", e)
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

	if ns.Sub != nil {
		if err := ns.Sub.Unsubscribe(); err != nil {
			ns.Logger.Infof("Unable to unsubscribe at %s: %s", ns.channel, err)
		}
	}

	if ns.conn.client != nil {
		ns.Logger.Infof("Closing connection: [%s: %s]", ns.conn.clientId, ns.channel)
		if err := ns.conn.client.Close(); err != nil {
			ns.Logger.Errorf("Unable to close nats streaming server connection at %s: %s", ns.channel, err)
		}
	}
}
