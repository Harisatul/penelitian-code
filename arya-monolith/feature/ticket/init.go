package ticket

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	db    *pgxpool.Pool
	cache *redis.Client
)

func SetDBPool(pool *pgxpool.Pool) {
	if pool == nil {
		panic("cannot assign nil db pool")
	}

	db = pool
}

func SetCachePool(pool *redis.Client) {
	if pool == nil {
		panic("cannot assign nil cache pool")
	}

	cache = pool

}
