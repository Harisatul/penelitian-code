package cmd

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
	"monolith/feature/shared"
	"monolith/feature/ticket"
	"time"
)

func runListTicketUpdater(ctx context.Context) {
	cfg := shared.LoadConfig("config/list_ticket_updater.yaml")

	dbCfg, err := pgxpool.ParseConfig(cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatalln("unable to parse database config", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Fatalln("unable to create database connection pool", err)
	}
	defer pool.Close()

	cacheClient := redis.NewClient(&redis.Options{
		Addr:                  cfg.RedisConfig.Addr,
		MinIdleConns:          int(cfg.RedisConfig.MinConn),
		ConnMaxIdleTime:       5 * time.Minute,
		ContextTimeoutEnabled: true,
	})
	defer cacheClient.Close()

	ticket.SetDBPool(pool)
	ticket.SetCachePool(cacheClient)

	ticker := time.NewTicker(3 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				ticket.UpdateListUnreservedTickets(ctx)
			}
		}
	}()

	log.Println("updater started")

	<-ctx.Done()
	ticker.Stop()

	log.Println("updater shutdown")
}
