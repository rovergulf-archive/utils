package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/rovergulf/utils/clog"
	"strings"
	"sync"
)

type Message struct {
	Offset    int64
	Partition int32
	Message   *sarama.ConsumerMessage
}

type Consumer struct {
	messages chan Message
	errors   chan error
	quit     chan struct{}
	topic    string
	client   sarama.Client
	consumer sarama.Consumer

	pom map[int32]sarama.PartitionOffsetManager
	wg  *sync.WaitGroup

	UseOldestOnFail bool
}

func NewConsumer(brokersAddr, topic, clientId string) (*Consumer, error) {

	var wg sync.WaitGroup

	c := Consumer{
		messages: make(chan Message),
		errors:   make(chan error),
		quit:     make(chan struct{}),
		topic:    topic,
		pom:      make(map[int32]sarama.PartitionOffsetManager),
		wg:       &wg,
	}

	var err error

	config := sarama.NewConfig()
	config.ClientID = clientId
	groupId := clientId

	config.Consumer.Offsets.Initial = sarama.OffsetOldest // configuration for LastOffset()

	addresses := strings.Split(brokersAddr, ",")
	if c.client, err = sarama.NewClient(addresses, config); err != nil {
		return nil, err
	}

	if c.consumer, err = sarama.NewConsumerFromClient(c.client); err != nil {
		return nil, err
	}

	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return nil, err
	}

	offsetManager, err := sarama.NewOffsetManagerFromClient(groupId, c.client)
	if err != nil {
		return nil, err
	}

	for _, p := range partitions {
		if c.pom[p], err = offsetManager.ManagePartition(c.topic, p); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c Consumer) Start(offsets map[int32]int64) error {

	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return fmt.Errorf("failed to retrieve numbers of partitions, %s", err)
	}

	for _, p := range partitions {
		offset, ok := offsets[p]
		if !ok {
			// here we give the newest-offset as a safeguard mechanism. If caller
			// by accident passed nil offsets, we won't start from the
			// beginning which could possibly write the same message multiple
			// times. If the caller always uses `LastOffset()` it will use the
			// correct offset.  Which is the oldest offset for the  first time
			// or continue from last commited offset.
			clog.Warningf("offset partition %d is not given, using newest offset", p)
			offset = sarama.OffsetNewest
		}

		clog.Infof("consuming kafka topic:%s partition:%d offset:%d", c.topic, p, offset)
		pc, err := c.consumer.ConsumePartition(c.topic, p, offset)
		if err != nil {
			if !c.UseOldestOnFail {
				return err
			}

			if strings.Contains(err.Error(), "The requested offset is outside the range") {
				pc, err = c.consumer.ConsumePartition(c.topic, p, sarama.OffsetOldest)
			}
		}

		messages := pc.Messages()
		errors := pc.Errors()
		go func(pc sarama.PartitionConsumer) {
			c.wg.Add(1)
			defer c.wg.Done()

			for {
				select {
				case <-c.quit:
					if pc != nil {
						err := pc.Close()
						if err != nil {
							clog.Errorf("Unable to close partition consumer")
						}
					}
					return
				case m := <-messages:
					c.messages <- Message{Offset: m.Offset, Partition: m.Partition, Message: m}
				case err := <-errors:
					c.errors <- err
				}
			}
		}(pc)
	}

	return nil
}

// stop consuming
func (c Consumer) Stop() {
	close(c.quit)
	clog.Info("Waiting kafka partition consumer to stop")
	c.wg.Wait()
	clog.Info("Kafka partition consumer stopped")
	if c.consumer != nil {
		err := c.consumer.Close()
		if err != nil {
			clog.Errorf("Unable to stop consumer")
		}
	}
	if c.pom != nil {
		for _, p := range c.pom {
			if p != nil {
				err := p.Close()
				if err != nil {
					clog.Errorf("Unable to stop partition offset manager")
				}
			}
		}
	}
	if c.client != nil {
		err := c.client.Close()
		if err != nil {
			clog.Errorf("Unable to stop partition consumer")
		}
	}
}

// Messages returns channel where new Message instances will be published.
func (c Consumer) Messages() <-chan Message {
	return c.messages
}

// Errors returns error while consuming kafka messages.
func (c Consumer) Errors() <-chan error {
	return c.errors
}

// LastOffset returns the last offset of the Client. This uses Kafka's offset
// manager API. If not found, it returns the oldest-offset for each partition
func (c Consumer) LastOffset() map[int32]int64 {
	offsets := make(map[int32]int64)
	for p, om := range c.pom {
		offsets[p], _ = om.NextOffset()
	}
	return offsets
}

// Ack offsets to be marked as processed on kafka offset-manager.
// Note: sarama MarkOffset is buffered (default 1s).
// Upon unclean shutdown we could double process the last 1s messages
func (c Consumer) Ack(partition int32, offset int64) {
	if om, ok := c.pom[partition]; ok {
		om.MarkOffset(offset, "")
	} else {
		clog.Warningf("Skipping Ack for unmanaged partition %d offset %d", partition, offset)
	}
}
