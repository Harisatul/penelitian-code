package cmd

import (
	"category-svc/feature/shared"
	"category-svc/feature/ticket"
	"category-svc/pkg"
	"context"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func runCreateOrderConsumer(ctx context.Context) {
	cfg := shared.LoadConfig("config/create_order_consumer.yaml")
	kafkaCfg := pkg.NewKafkaConsumerConfig()

	consumer, err := sarama.NewConsumerGroup([]string{cfg.Kafka.Broker}, ticket.CreateOrderConsumerGroup, kafkaCfg)
	if err != nil {
		log.Fatalln("unable to create consumer group", err)
	}

	defer consumer.Close()

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

	ticket.SetDBPool(pool)

	go func() {
		for err = range consumer.Errors() {
			log.Printf("consumer error, topic %s, error %s", ticket.CreateOrderTopic, err.Error())
		}
	}()

	go func() {
		for {
			select {
			case <-newCtx.Done():
				log.Println("consumer stopped")
				return
			default:
				err = consumer.Consume(newCtx, []string{ticket.CreateOrderTopic},
					pkg.NewKafkaConsumer(&ticket.CreateOrderHandler{}, 500),
				)
				if err != nil {
					log.Printf("consume message error, topic %s, error %s", ticket.CreateOrderTopic, err.Error())
					return
				}
			}
		}
	}()

	log.Printf("consumer up and running, topic %s, group: %s", ticket.CreateOrderTopic, ticket.CreateOrderConsumerGroup)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	cancel()
	log.Println("cancelled message without marking offsets")
}
