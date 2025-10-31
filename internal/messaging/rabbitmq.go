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

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(
		"update_request",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQService{
		Conn:    conn,
		Channel: channel,
	}, nil
}

func (s *RabbitMQService) Publish(message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return s.Channel.Publish(
		"",
		"update_request",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (s *RabbitMQService) Close() {
	s.Channel.Close()
	s.Conn.Close()
}
