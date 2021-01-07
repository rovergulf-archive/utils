package kafka

import (
	"context"
	"github.com/rovergulf/utils/clog"
)

type HighLevelConsumer struct {
	Consumer *Consumer
	Ctx      context.Context
}

func NewHighLevelConsumer(ctx context.Context, addr string, topic string, app string) *HighLevelConsumer {
	consumer, err := NewConsumer(addr, topic, app)
	if err != nil {
		clog.Fatal(err.Error())
	}
	if err := consumer.Start(consumer.LastOffset()); err != nil {
		clog.Fatal(err.Error())
	}
	return &HighLevelConsumer{
		Ctx:      ctx,
		Consumer: consumer,
	}
}

func (con *HighLevelConsumer) StartConsumption(handler func(value []byte) error) {
loop:
	for {
		select {
		case <-con.Ctx.Done():
			clog.Info("Received shutdown signal, stopping consumption")
			con.Consumer.Stop()
			break loop
		case m := <-con.Consumer.Messages():
			partition := m.Partition
			offset := m.Offset
			msg := m.Message
			if err := handler(msg.Value); err != nil {
				clog.Error(err.Error())
			}
			con.Consumer.Ack(partition, offset)
		case e := <-con.Consumer.Errors():
			clog.Error(e.Error())
		}
	}
}
