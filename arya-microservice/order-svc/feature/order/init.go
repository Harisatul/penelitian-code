package order

import (
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	db    *pgxpool.Pool
	kp    sarama.SyncProducer
	cache *redis.Client
)

func SetDBPool(dbPool *pgxpool.Pool) {
	if dbPool == nil {
		panic("cannot assign nil db pool")
	}

	db = dbPool
}

func SetKafkaProducer(producer sarama.SyncProducer) {
	if producer == nil {
		panic("cannot assign nil kafka producer")
	}

	kp = producer
}

func SetCachePool(cacheClient *redis.Client) {
	if cacheClient == nil {
		panic("cannot assign nil cache client")
	}

	cache = cacheClient
}
