package natsmq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"log"
	"time"
)

type StanConn struct {
	clientId string
	client   stan.Conn
	Tracer   opentracing.Tracer
	logger   *zap.SugaredLogger
}

//func NewStreamConn(c *Config, clientId string, opts ...stan.Option) (stan.Conn, error) {
//
//	// Send PINGs every 15 seconds, and fail after 5 PINGs without any response.
//	opts = append(opts, stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
//		log.Printf("! Nats streaming connection lost: %s", err)
//	}))
//	opts = append(opts, stan.Pings(15, 5))
//	opts = append(opts, stan.PubAckWait(stan.DefaultAckWait)) // 30 * time.Second
//
//	sc, err := stan.Connect(c.ClusterId, clientId, opts...)
//	if err != nil {
//		return nil, err
//	}
//
//	//log.Printf("[%s] Successfully connected to '%s' NATS-streaming cluster", clientId, clusterId)
//	return sc, nil
//}

func NewStanConn(lg *zap.SugaredLogger, c *Config, clientId string, opts ...stan.Option) (*StanConn, error) {
	s := new(StanConn)
	s.logger = lg
	s.clientId = fmt.Sprintf("%s-%d", clientId, time.Now().Unix())

	nopts := setupDefaultNatsConnOptions(lg, nil)
	nopts = append(nopts, nats.Name(clientId))

	nc, err := nats.Connect(c.Broker, nopts...)
	if err != nil {
		s.logger.Errorf("Failed to set nats server connection: %s", err)
		return nil, err
	}

	opts = append(opts, stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
		log.Printf("! Nats streaming connection lost: %s", err)
	}))
	opts = append(opts, stan.Pings(15, 5))
	opts = append(opts, stan.PubAckWait(stan.DefaultAckWait)) // 30 * time.Second
	opts = append(opts, stan.NatsConn(nc))

	sc, err := stan.Connect(c.ClusterId, clientId, opts...)
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

func DefaultAckHandler(nuid string, err error) {
	if err != nil {
		log.Printf("! Error publishing message [nuid: %s]: %s", nuid, err)
	} else {
		log.Printf("Received ack for message [nuid: %s]", nuid)
	}
}

func (sc *StanConn) SendMessage(topic string, data interface{}) {
	if sc.client == nil {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("Unable to marshal data: %s", err)
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
		log.Printf("Unable to marshal data: %s", err)
		return
	}

	res, err := sc.client.PublishAsync(topic, payload, handler)
	if err != nil {
		log.Printf("Error publishing message [%s: %s] due: %s", sc.clientId, topic, err)
	} else {
		log.Printf("Published to [%s]: [nuid: %s]", topic, res)
	}
}
