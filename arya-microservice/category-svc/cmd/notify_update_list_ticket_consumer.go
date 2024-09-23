package cmd

import (
	"category-svc/feature/shared"
	"category-svc/feature/ticket"
	"category-svc/pkg"
	"context"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func runNotifyUpdateListTicketConsumer(ctx context.Context) {
	cfg := shared.LoadConfig("config/notify_update_list_ticket_consumer.yaml")
	kafkaCfg := pkg.NewKafkaConsumerConfig()

	consumer, err := sarama.NewConsumerGroup([]string{cfg.Kafka.Broker}, ticket.UpdateListUnreservedTicketsConsumerGroup, kafkaCfg)
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

	cacheClient := redis.NewClient(&redis.Options{
		Addr:                  cfg.RedisConfig.Addr,
		MinIdleConns:          int(cfg.RedisConfig.MinConn),
		ConnMaxIdleTime:       5 * time.Minute,
		ContextTimeoutEnabled: true,
	})
	defer cacheClient.Close()

	ticket.SetDBPool(pool)
	ticket.SetCachePool(cacheClient)

	go func() {
		for err = range consumer.Errors() {
			log.Printf("consumer error, topic %s, error %s", ticket.NotifyUpdateListUnreservedTicketsTopic, err.Error())
		}
	}()

	go func() {
		for {
			select {
			case <-newCtx.Done():
				log.Println("consumer stopped")
				return
			default:
				err = consumer.Consume(newCtx, []string{ticket.NotifyUpdateListUnreservedTicketsTopic},
					pkg.NewKafkaConsumer(&ticket.UpdateListUnreservedTicketHandler{}, 1),
				)
				if err != nil {
					log.Printf("consume message error, topic %s, error %s", ticket.NotifyUpdateListUnreservedTicketsTopic, err.Error())
					return
				}
			}
		}
	}()

	log.Printf("consumer up and running, topic %s, group: %s", ticket.NotifyUpdateListUnreservedTicketsTopic, ticket.UpdateListUnreservedTicketsConsumerGroup)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	cancel()
	log.Println("cancelled message without marking offsets")
}
