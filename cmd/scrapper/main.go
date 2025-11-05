package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KOJIMEISTER/it_russian_stat/internal/domain"
	"github.com/KOJIMEISTER/it_russian_stat/internal/messaging"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	rabbitService, err := messaging.NewRabbitMQService()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitService.Close()

	msgs, err := rabbitService.Consume("update_request")
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	done := make(chan struct{})

	go processMessages(ctx, rabbitService, msgs, done)

	log.Printf(" [*] Waiting for messages. To exit precc CTRL+C")

	<-quit
	log.Println("Shutting down scrapper...")
	cancel()

	select {
	case <-done:
		log.Println("Scrapper shutdown gracefully")
	case <-time.After(10 * time.Second):
		log.Println("Scrapper shutdown timeout - forcing exit")
	}

}

func processMessages(ctx context.Context, rabbitService *messaging.RabbitMQService, msgs <-chan amqp091.Delivery, done chan<- struct{}) {
	defer close(done)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping message processing...")
			return
		case d, ok := <-msgs:
			if !ok {
				log.Println("Message channel closed")
				return
			}

			var req domain.UpdateRequestMessage
			msg := domain.ScrapperResponse{
				RequestID: req.RequestID,
				At:        time.Now().UTC(),
			}

			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Error decoding message: %v", err)
				msg.Status = "error"
				msg.StatusText = "wrong message format"
			} else {
				log.Printf("Recieved a message: %+v", req)
				msg.Status = "recieved"
				msg.StatusText = fmt.Sprintf("Scrapper accepted job for %s..%s", req.StartDate, req.EndDate)

				// go performBackgroundWork(ctx, rabbitService, req)
			}

			if err := rabbitService.PublishTo(req.CallbackQueue, msg); err != nil {
				log.Panicf("Failed to publish scrapper status: %v", err)
			}
		}
	}
}
