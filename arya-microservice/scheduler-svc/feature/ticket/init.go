package ticket

import (
	"github.com/IBM/sarama"
)

var (
	kp sarama.SyncProducer
)

func SetKafkaProducer(producer sarama.SyncProducer) {
	if producer == nil {
		panic("cannot assign nil kafka producer")
	}

	kp = producer
}
