package cmd

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"log"
	"os"
	"os/signal"
	"scheduler-svc/feature/order"
	"scheduler-svc/feature/shared"
	"scheduler-svc/pkg"
	"syscall"
	"time"
)

func runCreateOrderConsumer(ctx context.Context) {
	cfg := shared.LoadConfig("config/create_order_consumer.yaml")
	kafkaCfg := pkg.NewKafkaConsumerConfig()

	client, err := sarama.NewConsumerGroup([]string{cfg.Kafka.Broker}, order.CreateOrderConsumerGroup, kafkaCfg)
	if err != nil {
		log.Fatalln("unable to create consumer group", err)
	}

	defer client.Close()

	dbCfg, err := pgxpool.ParseConfig(cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatalln("unable to parse database config", err)
	}

	// Set needed dependencies
	newCtx, cancel := context.WithCancel(ctx)

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Fatalln("unable to create database connection pool", err)
	}
	defer pool.Close()

	queue, err := river.NewClient[pgx.Tx](riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		log.Fatalln("unable to create inserter job queue", err)
	}

	cacheClient := redis.NewClient(&redis.Options{
		Addr:            cfg.RedisConfig.Addr,
		MinIdleConns:    int(cfg.RedisConfig.MinConn),
		ConnMaxIdleTime: 5 * time.Minute,
	})

	order.SetQueue(queue)
	order.SetCachePool(cacheClient)

	go func() {
		for err = range client.Errors() {
			log.Printf("consumer error, topic %s, error %s", order.CreateOrderTopic, err.Error())
		}
	}()

	go func() {
		for {
			select {
			case <-newCtx.Done():
				log.Println("consumer stopped")
				return
			default:
				err = client.Consume(newCtx, []string{order.CreateOrderTopic},
					pkg.NewKafkaConsumer(&order.CreateOrderHandler{}, 1000),
				)
				if err != nil {
					log.Printf("consume message error, topic %s, error %s", order.CreateOrderTopic, err.Error())
					return
				}
			}
		}
	}()

	log.Printf("consumer up and running, topic %s, group: %s", order.CreateOrderTopic, order.CreateOrderConsumerGroup)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	cancel()
	log.Println("cancelled message without marking offsets")
}
