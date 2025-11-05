package config

import "os"

type RabbitMQConfig struct {
	URL string
}

type MongoDBConfig struct {
	URL string
}

func LoadRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		URL: os.Getenv("RABBITMQ_URI"),
	}
}

func LoadMongoDBConfig() MongoDBConfig {
	return MongoDBConfig{
		URL: os.Getenv("MONGODB_URI"),
	}
}
