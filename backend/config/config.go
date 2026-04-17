package config

import (
	"net/url"
	"time"
)

// Shared configuration interfaces

type DatabaseConfig interface {
	DBConnectionURL() string
}

type RedisConfig interface {
	RedisAddr() string
	GetRedisDB() int
}

type SessionConfig interface {
	GetSessionTTL() time.Duration
}

type CacheConfig interface {
	GetCacheTTL() time.Duration
}

type EmailConfig interface {
	SMTPAddress() string
	GetSMTPFrom() string
}

type WebpushConfig interface {
	VAPIDConfig() (publicKey, privateKey, subject string)
}

// noWebpushConfig is a null implementation for processes that don't need webpush.
type noWebpushConfig struct{}

func (n noWebpushConfig) VAPIDConfig() (string, string, string) {
	return "", "", ""
}

// NoWebpushConfig can be used when a process doesn't need to send push notifications.
var NoWebpushConfig WebpushConfig = noWebpushConfig{}

type QueueConfig interface {
	SQSConfig() (region, url, endpoint string)
}

type FrontendConfig interface {
	GetFrontendURL() *url.URL
}
