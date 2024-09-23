package pkg

import (
	"github.com/IBM/sarama"
	"time"
)

func NewKafkaProducerConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_6_0_0
	cfg.ChannelBufferSize = 1024
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Timeout = 3 * time.Second
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	return cfg
}

func PublishMessage(producer sarama.SyncProducer, topic, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}
	_, _, err := producer.SendMessage(msg)
	return err
}
