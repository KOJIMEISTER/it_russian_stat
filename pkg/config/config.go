package config

import "os"

type RabbitMQConfig struct {
	URL string
}

func LoadRabbitMQConfig() RabbitMQConfig {
	user := os.Getenv("RABBITMQ_USERNAME")
	pass := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	return RabbitMQConfig{
		URL: "amqp://" + user + ":" + pass + "@" + host + ":" + port + "/",
	}
}
