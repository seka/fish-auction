package config

import (
	"fmt"
	"net"
	"net/url"
	"time"
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
	SMTPHost         string
	SMTPPort         string
	SMTPFrom         string
	VAPIDPublicKey   string
	VAPIDPrivateKey  string
	VAPIDSubject     string
	FrontendURL      *url.URL
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
		SMTPHost:         GetEnv("SMTP_HOST", "mailhog"),
		SMTPPort:         GetEnv("SMTP_PORT", "1025"),
		SMTPFrom:         GetEnv("SMTP_FROM", "noreply@fish-auction.com"),
		VAPIDPublicKey:   GetEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey:  GetEnv("VAPID_PRIVATE_KEY", ""),
		VAPIDSubject:     GetEnv("VAPID_SUBJECT", "mailto:admin@example.com"),
		FrontendURL:      loadFrontendURL(),
		SQSQueueURL:      GetEnv("SQS_QUEUE_URL", "http://localhost:4566/000000000000/notification-queue"),
		SQSRegion:        GetEnv("SQS_REGION", "ap-northeast-1"),
		SQSEndpoint:      GetEnv("SQS_ENDPOINT", "http://localhost:4566"),
	}

	if err := validateBase(cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.FrontendURL); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *WorkerConfig) RedisAddr() string {
	return net.JoinHostPort(c.RedisHost, c.RedisPort)
}

func (c *WorkerConfig) GetRedisDB() int {
	return c.RedisDB
}

func (c *WorkerConfig) SMTPAddress() string {
	return fmt.Sprintf("%s:%s", c.SMTPHost, c.SMTPPort)
}

func (c *WorkerConfig) GetSMTPFrom() string {
	return c.SMTPFrom
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

func (c *WorkerConfig) GetFrontendURL() *url.URL {
	return c.FrontendURL
}

// Job server doesn't use session/cache TTL, but we provide defaults for interface satisfaction
func (c *WorkerConfig) GetSessionTTL() time.Duration {
	return 0
}

func (c *WorkerConfig) GetCacheTTL() time.Duration {
	return 0
}
