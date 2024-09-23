package pkg

import (
	"context"
	"github.com/IBM/sarama"
	"sync"
	"time"
)

func NewKafkaConsumerConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_6_0_0
	cfg.Consumer.Group.Session.Timeout = 8 * time.Second
	cfg.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	cfg.ChannelBufferSize = 1024
	return cfg
}

type KafkaConsumerHandler interface {
	Handle(ctx context.Context, msg *sarama.ConsumerMessage)
}

type KafkaConsumer struct {
	Handler KafkaConsumerHandler
	sem     chan struct{}
	wg      sync.WaitGroup
}

func NewKafkaConsumer(handler KafkaConsumerHandler, limit int32) *KafkaConsumer {
	return &KafkaConsumer{
		Handler: handler,
		sem:     make(chan struct{}, limit),
	}
}

func (c *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		c.sem <- struct{}{} // Acquire a token
		c.wg.Add(1)
		go func(msg *sarama.ConsumerMessage) {
			defer c.wg.Done()
			defer func() { <-c.sem }() // Release the token

			c.Handler.Handle(session.Context(), msg)
			session.MarkMessage(msg, "")
		}(msg)
	}

	c.wg.Wait()

	return nil
}

func (c *KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
