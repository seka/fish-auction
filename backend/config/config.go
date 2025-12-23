package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	ServerAddress   string
	RedisAddr       string
	CacheTTL        time.Duration
	AppEnv          string
	SMTPHost        string
	SMTPPort        string
	SMTPFrom        string
	VAPIDPublicKey  string
	VAPIDPrivateKey string
	VAPIDSubject    string
}

func Load() (*Config, error) {
	cacheTTL := getEnvInt("CACHE_TTL_SECONDS", 300) // デフォルト5分

	cfg := &Config{
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		ServerAddress:   os.Getenv("SERVER_ADDRESS"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		CacheTTL:        time.Duration(cacheTTL) * time.Second,
		AppEnv:          getEnv("APP_ENV", "production"),
		SMTPHost:        getEnv("SMTP_HOST", "mailhog"),
		SMTPPort:        getEnv("SMTP_PORT", "1025"),
		SMTPFrom:        getEnv("SMTP_FROM", "noreply@fish-auction.com"),
		VAPIDPublicKey:  getEnv("VAPID_PUBLIC_KEY", ""),
		VAPIDPrivateKey: getEnv("VAPID_PRIVATE_KEY", ""),
		VAPIDSubject:    getEnv("VAPID_SUBJECT", "mailto:admin@example.com"),
	}

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ":8080"
	}

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required environment variables")
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

func (c *Config) DBConnectionURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func (c *Config) SMTPAddress() string {
	return fmt.Sprintf("%s:%s", c.SMTPHost, c.SMTPPort)
}
