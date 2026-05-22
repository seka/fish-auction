package config

import (
	"fmt"
)

// RelayConfig represents the configuration for the outbox relay process.
type RelayConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSslMode  string
	AppEnv           string
	SQSQueueURL      string
	SQSRegion        string
	SQSEndpoint      string
}

// NewRelayConfig loads configuration for the relay process.
func NewRelayConfig() *RelayConfig {
	return &RelayConfig{
		PostgresHost:     GetEnv("POSTGRES_HOST", ""),
		PostgresPort:     GetEnv("POSTGRES_PORT", ""),
		PostgresUser:     GetEnv("POSTGRES_USER", ""),
		PostgresPassword: GetEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       GetEnv("POSTGRES_DB", ""),
		PostgresSslMode:  GetEnv("POSTGRES_SSLMODE", "disable"),
		AppEnv:           GetEnv("APP_ENV", "development"),
		SQSQueueURL:      GetEnv("AWS_SQS_QUEUE_URL", "http://localhost:4566/000000000000/notification-queue"),
		SQSRegion:        GetEnv("AWS_SQS_REGION", "ap-northeast-1"),
		SQSEndpoint:      GetEnv("AWS_SQS_ENDPOINT", "http://localhost:4566"),
	}
}

func (c *RelayConfig) DBConnectionURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresSslMode)
}

func (c *RelayConfig) SQSConfig() (region, queueURL, endpoint string) {
	return c.SQSRegion, c.SQSQueueURL, c.SQSEndpoint
}

func (c *RelayConfig) Validate() error {
	return validateSSLMode(c.AppEnv, c.PostgresSslMode)
}
