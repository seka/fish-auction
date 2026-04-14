package config

import (
	"fmt"
	"net"
	"net/url"
	"time"
)

// Config represents the application configuration loaded from environment variables.
type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	ServerHost       string
	ServerPort       string
	RedisHost        string
	RedisPort        string
	RedisDB          int
	CacheTTL         time.Duration
	SessionTTL       time.Duration
	AppEnv           string
	AllowedOrigins   string
	SMTPHost         string
	SMTPPort         string
	SMTPFrom         string
	VAPIDPublicKey   string
	VAPIDPrivateKey  string
	VAPIDSubject     string
	PostgresSslMode  string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	FrontendURL      *url.URL
	SQSQueueURL      string
	SQSRegion        string
	SQSEndpoint      string
}

// Load provides Load related functionality.
func Load() (*Config, error) {
	cacheTTL := GetEnvInt("CACHE_TTL_SECONDS", 300)       // デフォルト5分
	sessionTTL := GetEnvInt("SESSION_TTL_SECONDS", 86400) // デフォルト24時間

	cfg := &Config{
		PostgresHost:     GetEnv("POSTGRES_HOST", ""),
		PostgresPort:     GetEnv("POSTGRES_PORT", ""),
		PostgresUser:     GetEnv("POSTGRES_USER", ""),
		PostgresPassword: GetEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       GetEnv("POSTGRES_DB", ""),
		ServerHost:       GetEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:       GetEnv("SERVER_PORT", "8080"),
		RedisHost:        GetEnv("REDIS_HOST", "localhost"),
		RedisPort:        GetEnv("REDIS_PORT", "6379"),
		RedisDB:          GetEnvInt("REDIS_DB", 0),
		CacheTTL:         time.Duration(cacheTTL) * time.Second,
		SessionTTL:       time.Duration(sessionTTL) * time.Second,
		AppEnv:           GetEnv("APP_ENV", "develop"),
		AllowedOrigins:   GetEnv("ALLOWED_ORIGINS", "https://localhost,http://localhost:3000"),
		SMTPHost:         GetEnv("SMTP_HOST", "mailhog"),
		SMTPPort:         GetEnv("SMTP_PORT", "1025"),
		SMTPFrom:         GetEnv("SMTP_FROM", "noreply@fish-auction.com"),
		VAPIDPublicKey:   GetEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey:  GetEnv("VAPID_PRIVATE_KEY", ""),
		VAPIDSubject:     GetEnv("VAPID_SUBJECT", "mailto:admin@example.com"),
		PostgresSslMode:  GetEnv("POSTGRES_SSLMODE", "disable"),
		ReadTimeout:      time.Duration(GetEnvInt("SERVER_READ_TIMEOUT_SEC", 60)) * time.Second,
		WriteTimeout:     time.Duration(GetEnvInt("SERVER_WRITE_TIMEOUT_SEC", 60)) * time.Second,
		IdleTimeout:      time.Duration(GetEnvInt("SERVER_IDLE_TIMEOUT_SEC", 60)) * time.Second,
		FrontendURL: func() *url.URL {
			frontendURLStr := GetEnv("FRONTEND_URL", "https://localhost")
			frontendURL, err := url.Parse(frontendURLStr)
			if err != nil {
				return nil
			}
			if frontendURL.Scheme == "" || frontendURL.Host == "" {
				return nil
			}
			return frontendURL
		}(),
		SQSQueueURL:      getEnv("SQS_QUEUE_URL", "http://localhost:4566/000000000000/notification-queue"),
		SQSRegion:        getEnv("SQS_REGION", "ap-northeast-1"),
		SQSEndpoint:      getEnv("SQS_ENDPOINT", "http://localhost:4566"),
	}

	if cfg.PostgresHost == "" || cfg.PostgresPort == "" || cfg.PostgresUser == "" || cfg.PostgresPassword == "" || cfg.PostgresDB == "" {
		return nil, fmt.Errorf("missing required environment variables (POSTGRES_*)")
	}

	if cfg.FrontendURL == nil {
		return nil, fmt.Errorf("invalid or missing FRONTEND_URL")
	}

	return cfg, nil
}

// ServerAddr returns the joined host and port address for the server.
func (c *Config) ServerAddr() string {
	return net.JoinHostPort(c.ServerHost, c.ServerPort)
}

// RedisAddr returns the joined host and port address for the Redis server.
func (c *Config) RedisAddr() string {
	return net.JoinHostPort(c.RedisHost, c.RedisPort)
}

// DBConnectionURL returns the PostgreSQL connection string based on the configuration.
func (c *Config) DBConnectionURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresSslMode)
}

// SMTPAddress returns the SMTP server address in host:port format.
func (c *Config) SMTPAddress() string {
	return fmt.Sprintf("%s:%s", c.SMTPHost, c.SMTPPort)
}
