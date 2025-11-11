package config

import "os"

type RabbitMQConfig struct {
	URL string
}

type MongoDBConfig struct {
	URL        string
	Database   string
	Collection string
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

func LoadMongoDBConfig() MongoDBConfig {
	host := os.Getenv("MONGODB_HOST")
	user := os.Getenv("MONGODB_USERNAME")
	port := os.Getenv("MONGODB_PORT")
	database := os.Getenv("MONGODB_DATABASE")
	pass := os.Getenv("MONGODB_PASSWORD")
	return MongoDBConfig{
		URL:        "mongodb://" + user + ":" + pass + "@" + host + ":" + port + "/" + database + "?authSource=" + user,
		Database:   database,
		Collection: "vacancy",
	}
}
