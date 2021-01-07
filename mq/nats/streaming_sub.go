package natsmq

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/opentracing/opentracing-go"
	"log"
	"path"
	"strings"
	"time"
)

type StanSub struct {
	ctx        context.Context
	Tracer     opentracing.Tracer
	conn       *StanConn
	sub        stan.Subscription
	messages   chan *stan.Msg
	errors     chan error
	quit       chan struct{}
	channel    string
	ackManager *AckManager
}

func NewChanSubWithTracer(ctx context.Context, tracer opentracing.Tracer, storageDirPath, clusterId, channel, brokerAddr string, opts ...nats.Option) (*StanSub, error) {
	ns, err := NewChanSub(ctx, storageDirPath, clusterId, channel, brokerAddr, opts...)
	if err != nil {
		return nil, err
	}

	ns.Tracer = tracer
	return ns, nil
}

// NewChanSub creates connection with channel-named clientID
// creating subscription with a whole service lifetime context
func NewChanSub(ctx context.Context, storageDirPath, clusterId, brokerAddr, channel string, opts ...nats.Option) (*StanSub, error) {
	s := new(StanSub)
	s.ctx = ctx
	s.messages = make(chan *stan.Msg)
	s.errors = make(chan error)
	s.quit = make(chan struct{})
	s.channel = channel

	var clientId string
	dotIndex := strings.Index(channel, ".")
	if dotIndex > 0 {
		clientId = fmt.Sprintf("%s-%s-chan", channel[:dotIndex], channel[dotIndex+1:])
	} else {
		clientId = fmt.Sprintf("%s-chan", channel)
	}

	flushInterval := 30 * time.Minute
	dumpPath := path.Join(storageDirPath, fmt.Sprintf("/%s.dump", clientId))
	log.Printf("[%s] Start execution NATS ack manager dump from %s with %s interval", clientId, dumpPath, flushInterval)
	s.ackManager = NewAckTimestampManager(dumpPath, flushInterval)
	lastSequence := s.ackManager.Get()

	conn, err := NewStanConn(ctx, clusterId, brokerAddr, clientId, opts...)
	if err != nil {
		log.Printf("Unable to connect [%s:%s]: %s", clientId, brokerAddr, err)
		return nil, err
	}
	s.conn = conn

	// set subscription options for fault tolerance
	// only ack manually
	var sopts []stan.SubscriptionOption
	sopts = append(sopts, stan.SetManualAckMode())
	sopts = append(sopts, stan.AckWait(60*time.Second))
	sopts = append(sopts, stan.StartAtSequence(lastSequence))

	sub, err := s.conn.client.Subscribe(channel, func(msg *stan.Msg) {
		s.messages <- msg
	}, sopts...)
	if err != nil {
		log.Printf("Unable to subscribe [%s: %s]: %s", clientId, channel, err)
		return nil, err
	}
	s.sub = sub

	log.Printf("[%s] NATS subscription started awating '%s'-channel at sequence %d", clientId, channel, lastSequence)
	return s, nil
}

func (s *StanSub) StartConsumption(ctx context.Context, handler func(data []byte) error) {
loop:
	for {
		select {
		case <-s.ctx.Done():
			log.Printf("[%s] Received shutdown signal, stopping '%s' subscription", s.conn.clientId, s.channel)
			s.Stop()
			break loop
		case msg := <-s.messages:
			var span opentracing.Span

			if s.Tracer != nil {
				span = s.Tracer.StartSpan(fmt.Sprintf("[%s:%d]", s.channel, msg.Sequence))
			}
			// check if we have available handler
			if msg.Sequence < s.ackManager.sequence {
				continue loop
			}

			if err := handler(msg.Data); err != nil {
				log.Printf("Unable to handle nats '%s' subscription message: %s", s.channel, err)
				s.errors <- err
			} else {
				log.Printf("[%s] Successfully received sequenced message: %d at %d", s.channel, msg.Sequence, msg.Timestamp)
			}

			if err := msg.Ack(); err != nil {
				log.Printf("Unable to respond nats message: %s", err)
			} else {
				s.ackManager.Set(msg.Sequence)
				log.Printf("[%s] Succesfully acked: %s", s.channel, msg.Reply)
			}

			if span != nil {
				span.Finish()
			}
		case e := <-s.errors:
			log.Printf("Subscription error: %s", e)
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
	if s.ackManager != nil {
		s.ackManager.Flush()
	}

	if s.sub != nil {
		if err := s.sub.Unsubscribe(); err != nil {
			log.Printf("Unable to unsubscribe at %s: %s", s.channel, err)
		}
	}

	if s.conn.client != nil {
		log.Printf("Closing connection: [%s: %s]", s.conn.clientId, s.channel)
		if err := s.conn.client.Close(); err != nil {
			log.Printf("Unable to close nats streaming server connection at %s: %s", s.channel, err)
		}
	}
}
