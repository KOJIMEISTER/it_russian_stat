package config

import "time"

type ClientConfig interface {
	GetBaseURL() string
	GetSearchUrl() string
	GetVacancyUrl() string
	GetBearerToken() string
	GetTimeout() time.Duration
}
