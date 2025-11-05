package config

import (
	"os"
	"strconv"
	"time"
)

type ScrapperConfig struct {
	RabbitMQ    RabbitMQConfig
	MongoDB     MongoDBConfig
	ExternalAPI ExternalAPIConfig
	Settings    ScrapperSettings
}

type ExternalAPIConfig struct {
	BaseURL     string
	BearerToken string
	Timeout     time.Duration
}

type ScrapperSettings struct {
	WorkerCount  int
	RequestDelay time.Duration
	RetryCount   int
}

func LoadScrapperConfig() *ScrapperConfig {
	return &ScrapperConfig{
		RabbitMQ: LoadRabbitMQConfig(),
		MongoDB:  LoadMongoDBConfig(),
		ExternalAPI: ExternalAPIConfig{
			BaseURL:     os.Getenv("HH_API_URL"),
			BearerToken: os.Getenv("HH_API_TOKEN"),
			Timeout:     10 * time.Second,
		},
		Settings: ScrapperSettings{
			WorkerCount:  getEnvAsInt("SCRAPPER_WORKER_COUNT", 12),
			RequestDelay: 10 * time.Second,
			RetryCount:   3,
		},
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
