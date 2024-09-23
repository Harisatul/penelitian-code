package order

import (
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/riverqueue/river"
)

var (
	queue *river.Client[pgx.Tx]
	kp    sarama.SyncProducer
	cache *redis.Client
)

func SetQueue(jobQueue *river.Client[pgx.Tx]) {
	if jobQueue == nil {
		panic("cannot assign nil job queue")
	}

	queue = jobQueue
}

func SetKafkaProducer(producer sarama.SyncProducer) {
	if producer == nil {
		panic("cannot assign nil kafka producer")
	}

	kp = producer
}

func SetCachePool(cachePool *redis.Client) {
	if cachePool == nil {
		panic("cannot assign nil cache pool")
	}

	cache = cachePool
}
