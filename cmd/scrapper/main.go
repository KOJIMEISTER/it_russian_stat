package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/KOJIMEISTER/it_russian_stat/internal/domain"
	"github.com/KOJIMEISTER/it_russian_stat/internal/messaging"
)

func main() {
	rabbitService, err := messaging.NewRabbitMQService()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitService.Close()

	msgs, err := rabbitService.Consume("update_request")
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
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
			}
			if err := rabbitService.PublishTo(req.CallbackQueue, msg); err != nil {
				log.Printf("Failed to publish scrapper status: %v", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit precc CTRL+C")
	<-forever
}
