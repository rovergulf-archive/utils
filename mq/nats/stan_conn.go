package natsmq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nuid"
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
	nuid     *nuid.NUID
}

func NewStanConn(c *Config) (*StanConn, error) {
	s := &StanConn{
		Tracer: c.Tracer,
		logger: c.Logger.Named("nats-streaming"),
		nuid:   nuid.New(),
	}
	s.clientId = fmt.Sprintf("%s-%d", c.ClientId, time.Now().Unix())

	nc, err := NewConn(c)
	if err != nil {
		s.logger.Errorw("Failed to set nats server connection",
			"broker", c.Broker, "err", err)
		return nil, err
	}

	c.StanConn = append(c.StanConn, stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
		s.logger.Warnf("Connection lost: %s", err)
	}))
	c.StanConn = append(c.StanConn, stan.Pings(15, 5))
	c.StanConn = append(c.StanConn, stan.NatsConn(nc))
	if c.AckTimeout > 0 {
		c.StanConn = append(c.StanConn, stan.PubAckWait(time.Duration(c.AckTimeout)))
	} else {
		c.StanConn = append(c.StanConn, stan.PubAckWait(stan.DefaultAckWait)) // 30 * time.Second
	}

	sc, err := stan.Connect(c.ClusterId, c.ClientId, c.StanConn...)
	if err != nil {
		s.logger.Errorw("Failed to set stan connection",
			"client_id", c.ClientId, "cluster_id", c.ClusterId, "err", err)
		return nil, err
	}
	s.client = sc

	s.logger.Infow("Initialized nats-streaming conn",
		"broker", c.Broker, "cluster_id", c.ClusterId, "client_id", c.ClientId)
	return s, nil
}

func (sc *StanConn) Stop() {
	if sc.client != nil {
		sc.logger.Debugf("Closing '%s' connection", sc.clientId)
		if err := sc.client.Close(); err != nil {
			sc.logger.Errorf("Unable to stop nats streaming server connection: %s", err)
		}
	}
}

func (sc *StanConn) DefaultAckHandler(gNuid string, err error) {
	if err != nil {
		sc.logger.Errorw("Error publishing message", "guid", gNuid, "err", err)
	} else {
		sc.logger.Infow("Received ack for message", "nuid", gNuid, "guid", sc.nuid.Next())
	}
}

func (sc *StanConn) SendMessage(channel string, data interface{}) {
	if sc.client == nil {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		sc.logger.Errorf("Unable to marshal data: %s", err)
		return
	}

	if err := sc.client.Publish(channel, payload); err != nil {
		sc.logger.Errorw("Error publishing message",
			"client_id", sc.clientId, "chan", channel, "err", err)
	} else {
		sc.logger.Infow("Published",
			"chan", channel, "guid", sc.nuid.Next(), "async", false)
	}
}

func (sc *StanConn) SendAsyncMessage(channel string, data interface{}) {
	if sc.client == nil {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		sc.logger.Errorf("Unable to marshal data: %s", err)
		return
	}

	res, err := sc.client.PublishAsync(channel, payload, sc.DefaultAckHandler)
	if err != nil {
		sc.logger.Errorw("Error publishing",
			"chan", channel, "client_id", sc.clientId, "err", err)
	} else {
		sc.logger.Infow("Published",
			"chan", channel, "nuid", res, "guid", sc.nuid.Next(), "async", true)
	}
}
