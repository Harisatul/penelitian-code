package cmd

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"log"
	"scheduler-svc/feature/order"
	"scheduler-svc/feature/shared"
	"scheduler-svc/pkg"
	"time"
)

func runScheduler(ctx context.Context) {
	cfg := shared.LoadConfig("config/scheduler.yaml")

	dbCfg, err := pgxpool.ParseConfig(cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatalln("unable to parse database config", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Fatalln("unable to create database connection pool", err)
	}
	defer pool.Close()

	producer, err := sarama.NewSyncProducer([]string{cfg.Kafka.Broker}, pkg.NewKafkaProducerConfig())
	if err != nil {
		log.Fatalln("unable to create kafka producer", err)
	}
	defer producer.Close()

	order.SetKafkaProducer(producer)

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
