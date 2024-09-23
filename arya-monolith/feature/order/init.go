package order

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/riverqueue/river"
)

var (
	db    *pgxpool.Pool
	queue *river.Client[pgx.Tx]
	cache *redis.Client
)

func SetDBPool(dbPool *pgxpool.Pool) {
	if dbPool == nil {
		panic("cannot assign nil db pool")
	}

	db = dbPool
}

func SetQueue(jobQueue *river.Client[pgx.Tx]) {
	if jobQueue == nil {
		panic("cannot assign nil job queue")
	}

	queue = jobQueue
}

func SetCachePool(cacheClient *redis.Client) {
	if cacheClient == nil {
		panic("cannot assign nil cache client")
	}

	cache = cacheClient
}
