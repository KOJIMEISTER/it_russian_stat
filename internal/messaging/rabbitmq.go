package messaging

import (
	"encoding/json"

	"github.com/KOJIMEISTER/it_russian_stat/pkg/config"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func NewRabbitMQService() (*RabbitMQService, error) {
	cfg := config.LoadRabbitMQConfig()

	conn, err := amqp091.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if _, err = ch.QueueDeclare("update_request", true, false, false, false, nil); err != nil {
		return nil, err
	}

	if _, err = ch.QueueDeclare("scrapper_status", true, false, false, false, nil); err != nil {
		return nil, err
	}

	return &RabbitMQService{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (s *RabbitMQService) PublishTo(queue string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return s.Channel.Publish(
		"",
		queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (s *RabbitMQService) Consume(queue string) (<-chan amqp091.Delivery, error) {
	return s.Channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (s *RabbitMQService) Close() {
	s.Channel.Close()
	s.Conn.Close()
}
