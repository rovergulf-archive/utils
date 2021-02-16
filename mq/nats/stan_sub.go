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
	sub      stan.Subscription
	messages chan *stan.Msg
	errors   chan error
	quit     chan struct{}
	channel  string
}

func NewChanSubWithTracer(lg *zap.SugaredLogger, t opentracing.Tracer, c *Config, channel string, opts ...stan.Option) (*StanSub, error) {
	ns, err := NewChanSub(lg, c, channel, opts...)
	if err != nil {
		return nil, err
	}

	ns.Tracer = t
	return ns, nil
}

// NewChanSub creates connection with channel-named clientID
// creating subscription with a whole service lifetime context
func NewChanSub(lg *zap.SugaredLogger, c *Config, channel string, sopts ...stan.Option) (*StanSub, error) {
	s := new(StanSub)
	s.Logger = lg.Named("stan-sub-" + channel)
	s.messages = make(chan *stan.Msg)
	s.errors = make(chan error)
	s.quit = make(chan struct{})
	s.channel = channel

	ch := strings.Split(channel, ",")
	clientId := fmt.Sprintf("%s-chan-%d", strings.Join(ch, "-"), time.Now().Unix())

	conn, err := NewStanConn(s.Logger, c, clientId, sopts...)
	if err != nil {
		s.Logger.Errorf("Unable to connect [%s:%s]: %s", clientId, c.Broker, err)
		return nil, err
	}
	s.conn = conn

	// set subscription options for fault tolerance
	// only ack manually
	var opts []stan.SubscriptionOption
	opts = append(opts, stan.SetManualAckMode())
	opts = append(opts, stan.AckWait(60*time.Second))
	opts = append(opts, stan.StartWithLastReceived())

	sub, err := s.conn.client.Subscribe(channel, func(msg *stan.Msg) {
		s.messages <- msg
	}, opts...)
	if err != nil {
		s.Logger.Errorf("Unable to subscribe [%s: %s]: %s", clientId, channel, err)
		return nil, err
	}
	s.sub = sub

	delivered, err := sub.Delivered()
	if err != nil {
		s.Logger.Errorf("Unable to check subscrpition delivered count: %s", err)
		return nil, err
	}

	s.Logger.Infof("[%s] NATS subscription started for '%s' awating at delivered count %d", clientId, channel, delivered)
	return s, nil
}

type StanSubHandler func(data []byte, sequence uint64, reply string) error

func (s *StanSub) StartConsumption(ctx context.Context, handler StanSubHandler) {
loop:
	for {
		select {
		case <-ctx.Done():
			s.Logger.Infof("[%s] Received shutdown signal, stopping '%s' subscription", s.conn.clientId, s.channel)
			s.Stop()
			break loop
		case msg := <-s.messages:
			var span opentracing.Span

			if s.Tracer != nil {
				span = s.Tracer.StartSpan(s.channel)
				span.SetTag("sequence", msg.Sequence)
				ctx = opentracing.ContextWithSpan(ctx, span)
			}

			if err := handler(msg.Data, msg.Sequence, msg.Reply); err != nil {
				s.Logger.Errorf("Unable to handle nats '%s' subscription message: %s", s.channel, err)
				s.errors <- err
			} else {
				s.Logger.Infof("[%s] Successfully received sequenced message: %d at %d", s.channel, msg.Sequence, msg.Timestamp)
			}

			if err := msg.Ack(); err != nil {
				s.Logger.Infof("Unable to respond nats message: %s", err)
			} else {
				s.Logger.Infow("Successfully acked",
					"channel", s.channel, "sequence", msg.Sequence, "reply", msg.Reply)
			}

			if span != nil {
				span.Finish()
			}
		case e := <-s.errors:
			s.Logger.Errorf("Subscription error: %s", e)
		}
	}
}

func (s *StanSub) Messages() <-chan *stan.Msg {
	return s.messages
}

func (s *StanSub) Errors() <-chan error {
	return s.errors
}

func (s *StanSub) Stop() {

	if s.sub != nil {
		if err := s.sub.Unsubscribe(); err != nil {
			s.Logger.Infof("Unable to unsubscribe at %s: %s", s.channel, err)
		}
	}

	if s.conn.client != nil {
		s.Logger.Infof("Closing connection: [%s: %s]", s.conn.clientId, s.channel)
		if err := s.conn.client.Close(); err != nil {
			s.Logger.Errorf("Unable to close nats streaming server connection at %s: %s", s.channel, err)
		}
	}
}
