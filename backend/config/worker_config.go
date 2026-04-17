package config

import (
	"fmt"
	"net"
)

// WorkerConfig represents the configuration for the background worker.
type WorkerConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSslMode  string
	RedisHost        string
	RedisPort        string
	RedisDB          int
	AppEnv           string
	VAPIDPublicKey   string
	VAPIDPrivateKey  string
	VAPIDSubject     string
	SQSQueueURL      string
	SQSRegion        string
	SQSEndpoint      string
}

// LoadWorkerConfig loads configuration for the background worker.
func LoadWorkerConfig() (*WorkerConfig, error) {
	cfg := &WorkerConfig{
		PostgresHost:     GetEnv("POSTGRES_HOST", ""),
		PostgresPort:     GetEnv("POSTGRES_PORT", ""),
		PostgresUser:     GetEnv("POSTGRES_USER", ""),
		PostgresPassword: GetEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       GetEnv("POSTGRES_DB", ""),
		PostgresSslMode:  GetEnv("POSTGRES_SSLMODE", "disable"),
		RedisHost:        GetEnv("REDIS_HOST", "localhost"),
		RedisPort:        GetEnv("REDIS_PORT", "6379"),
		RedisDB:          GetEnvInt("REDIS_DB", 0),
		AppEnv:           GetEnv("APP_ENV", "develop"),
		VAPIDPublicKey:   GetEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey:  GetEnv("VAPID_PRIVATE_KEY", ""),
		VAPIDSubject:     GetEnv("VAPID_SUBJECT", "mailto:admin@example.com"),
		SQSQueueURL:      GetEnv("SQS_QUEUE_URL", "http://localhost:4566/000000000000/notification-queue"),
		SQSRegion:        GetEnv("SQS_REGION", "ap-northeast-1"),
		SQSEndpoint:      GetEnv("SQS_ENDPOINT", "http://localhost:4566"),
	}

	return cfg, nil
}

func (c *WorkerConfig) RedisAddr() string {
	return net.JoinHostPort(c.RedisHost, c.RedisPort)
}

func (c *WorkerConfig) GetRedisDB() int {
	return c.RedisDB
}

func (c *WorkerConfig) VAPIDConfig() (string, string, string) {
	return c.VAPIDPublicKey, c.VAPIDPrivateKey, c.VAPIDSubject
}

func (c *WorkerConfig) SQSConfig() (string, string, string) {
	return c.SQSRegion, c.SQSQueueURL, c.SQSEndpoint
}

func (c *WorkerConfig) DBConnectionURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresSslMode)
}
