package natsmq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"time"
)

type StanConn struct {
	clientId string
	client   stan.Conn
	Tracer   opentracing.Tracer
	logger   *zap.SugaredLogger
}

func NewStanConnWithTracer(c *Config, tracer opentracing.Tracer) (*StanConn, error) {
	sc, err := NewStanConn(c)
	if err != nil {
		return nil, err
	}
	sc.Tracer = tracer

	return sc, nil
}

func NewStanConn(c *Config) (*StanConn, error) {
	s := new(StanConn)
	s.logger = c.Logger.Named("stan")
	s.clientId = fmt.Sprintf("%s-%d", c.ClientId, time.Now().Unix())

	nc, err := NewConn(c)
	if err != nil {
		s.logger.Errorf("Failed to set nats server connection: %s", err)
		return nil, err
	}

	c.StanConn = append(c.StanConn, stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
		s.logger.Warnf("Connection lost: %s", err)
	}))
	c.StanConn = append(c.StanConn, stan.Pings(15, 5))
	c.StanConn = append(c.StanConn, stan.PubAckWait(stan.DefaultAckWait)) // 30 * time.Second
	c.StanConn = append(c.StanConn, stan.NatsConn(nc))

	sc, err := stan.Connect(c.ClusterId, c.ClientId, c.StanConn...)
	if err != nil {
		s.logger.Errorf("Failed to set stan connection: %s", err)
		return nil, err
	}
	s.client = sc

	return s, nil
}

func (sc *StanConn) Stop() {
	if sc.client != nil {
		if err := sc.client.Close(); err != nil {
			sc.logger.Errorf("Unable to stop nats streaming server connection: %s", err)
		}
	}
}

func (sc *StanConn) DefaultAckHandler(nuid string, err error) {
	if err != nil {
		sc.logger.Errorf("Error publishing message [nuid: %s]: %s", nuid, err)
	} else {
		sc.logger.Infof("Received ack for message [nuid: %s]", nuid)
	}
}

func (sc *StanConn) SendMessage(topic string, data interface{}) {
	if sc.client == nil {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		sc.logger.Errorf("Unable to marshal data: %s", err)
		return
	}

	if err := sc.client.Publish(topic, payload); err != nil {
		sc.logger.Errorf("Error publishing message [%s: %s] due: %s", sc.clientId, topic, err)
	} else {
		sc.logger.Infof("Published message to [%s]", topic)
	}
}

func (sc *StanConn) SendAsyncMessage(topic string, data interface{}, handler stan.AckHandler) {
	if sc.client == nil {
		return
	}
	if handler == nil {
		handler = sc.DefaultAckHandler
	}

	payload, err := json.Marshal(data)
	if err != nil {
		sc.logger.Errorf("Unable to marshal data: %s", err)
		return
	}

	res, err := sc.client.PublishAsync(topic, payload, handler)
	if err != nil {
		sc.logger.Errorf("Error publishing message [%s: %s] due: %s", sc.clientId, topic, err)
	} else {
		sc.logger.Infof("Published to [%s]: [nuid: %s]", topic, res)
	}
}
