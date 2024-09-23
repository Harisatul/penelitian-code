package cmd

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"log"
	"monolith/feature/order"
	"monolith/feature/shared"
	"time"
)

func runJobQueue(ctx context.Context) {
	// Load configuration
	cfg := shared.LoadConfig("config/job_queue.yaml")

	dbCfg, err := pgxpool.ParseConfig(cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatalln("unable to parse database config", err)
	}

	// Set needed dependencies
	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Fatalln("unable to create database connection pool", err)
	}
	defer pool.Close()

	order.SetDBPool(pool)

	workers := river.NewWorkers()
	err = river.AddWorkerSafely[order.CancellationArgs](workers, &order.CancellationWorker{})
	if err != nil {
		log.Fatalln("unable to add worker", err)
	}

	queue, err := river.NewClient[pgx.Tx](riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 1000},
		},
		Workers: workers,
	})
	if err != nil {
		log.Fatalln("unable to create job queue", err)
	}

	go func() {
		if err := queue.Start(ctx); err != nil {
			log.Fatalln("unable to start job queue", err)
		}
	}()

	log.Println("job queue started")

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := queue.Stop(ctxShutDown); err != nil {
		log.Fatalln("unable to stop job queue", err)
	}

	log.Println("job queue stopped")
}
