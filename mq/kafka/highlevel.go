package kafka

import (
	"context"
	"go.uber.org/zap"
)

type HighLevelConsumer struct {
	Consumer *Consumer
	Ctx      context.Context
}

func NewHighLevelConsumer(ctx context.Context, lg *zap.SugaredLogger, addr string, topic string, app string) *HighLevelConsumer {
	consumer, err := NewConsumer(lg, addr, topic, app)
	if err != nil {
		lg.Fatal(err)
	}
	if err := consumer.Start(consumer.LastOffset()); err != nil {
		lg.Fatal(err.Error())
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
			con.Consumer.Logger.Info("Received shutdown signal, stopping consumption")
			con.Consumer.Stop()
			break loop
		case m := <-con.Consumer.Messages():
			partition := m.Partition
			offset := m.Offset
			msg := m.Message
			if err := handler(msg.Value); err != nil {
				con.Consumer.Logger.Error(err)
			}
			con.Consumer.Ack(partition, offset)
		case e := <-con.Consumer.Errors():
			con.Consumer.Logger.Error(e)
		}
	}
}
