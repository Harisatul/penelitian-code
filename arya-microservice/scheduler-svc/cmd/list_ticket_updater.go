package cmd

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"scheduler-svc/feature/shared"
	"scheduler-svc/feature/ticket"
	"scheduler-svc/pkg"
	"time"
)

func runListTicketUpdater(ctx context.Context) {
	cfg := shared.LoadConfig("config/list_ticket_updater.yaml")

	producer, err := sarama.NewSyncProducer([]string{cfg.Kafka.Broker}, pkg.NewKafkaProducerConfig())
	if err != nil {
		log.Fatalln("unable to create kafka producer", err)
	}
	defer producer.Close()

	ticket.SetKafkaProducer(producer)

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
