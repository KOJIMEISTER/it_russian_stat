package main

import (
	"encoding/json"
	"log"

	"github.com/KOJIMEISTER/it_russian_stat/internal/domain"
	"github.com/KOJIMEISTER/it_russian_stat/internal/messaging"
)

func main() {
	rabbitService, err := messaging.NewRabbitMQService()
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitService.Close()

	msgs, err := rabbitService.Channel.Consume(
		"update_request",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var req domain.UpdateRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			log.Printf("Recieved a message: %+v", req)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit precc CTRL+C")
	<-forever
}
