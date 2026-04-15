package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
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
	RedisAddr        string
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
}

// Load provides Load related functionality.
func Load() (*Config, error) {
	cacheTTL := getEnvInt("CACHE_TTL_SECONDS", 300)       // デフォルト5分
	sessionTTL := getEnvInt("SESSION_TTL_SECONDS", 86400) // デフォルト24時間

	cfg := &Config{
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		ServerHost:       os.Getenv("SERVER_HOST"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:          getEnvInt("REDIS_DB", 0),
		CacheTTL:         time.Duration(cacheTTL) * time.Second,
		SessionTTL:       time.Duration(sessionTTL) * time.Second,
		AppEnv:           getEnv("APP_ENV", "production"),
		AllowedOrigins:   getEnv("ALLOWED_ORIGINS", "https://localhost,http://localhost:3000"),
		SMTPHost:         getEnv("SMTP_HOST", "mailhog"),
		SMTPPort:         getEnv("SMTP_PORT", "1025"),
		SMTPFrom:         getEnv("SMTP_FROM", "noreply@fish-auction.com"),
		VAPIDPublicKey:   getEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey:  getEnv("VAPID_PRIVATE_KEY", ""),
		VAPIDSubject:     getEnv("VAPID_SUBJECT", "mailto:admin@example.com"),
		PostgresSslMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		ReadTimeout:      time.Duration(getEnvInt("SERVER_READ_TIMEOUT_SEC", 60)) * time.Second,
		WriteTimeout:     time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT_SEC", 60)) * time.Second,
		IdleTimeout:      time.Duration(getEnvInt("SERVER_IDLE_TIMEOUT_SEC", 60)) * time.Second,
		FrontendURL: func() *url.URL {
			frontendURLStr := getEnv("FRONTEND_URL", "https://localhost")
			frontendURL, err := url.Parse(frontendURLStr)
			if err != nil {
				return nil
			}
			if frontendURL.Scheme == "" || frontendURL.Host == "" {
				return nil
			}
			return frontendURL
		}(),
	}

	if cfg.PostgresHost == "" || cfg.PostgresPort == "" || cfg.PostgresUser == "" || cfg.PostgresPassword == "" || cfg.PostgresDB == "" {
		return nil, fmt.Errorf("missing required environment variables (POSTGRES_*)")
	}

	if cfg.FrontendURL == nil {
		return nil, fmt.Errorf("invalid or missing FRONTEND_URL")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// ServerAddr returns the joined host and port address for the server.
func (c *Config) ServerAddr() string {
	return net.JoinHostPort(c.ServerHost, c.ServerPort)
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
