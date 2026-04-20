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

// noSessionConfig is a null implementation for processes that don't need sessions.
type noSessionConfig struct{}

func (n noSessionConfig) GetSessionTTL() time.Duration { return 0 }

// NoSessionConfig can be used when a process doesn't need to handle sessions.
var NoSessionConfig SessionConfig = noSessionConfig{}

type CacheConfig interface {
	GetCacheTTL() time.Duration
}

// noCacheConfig is a null implementation for processes that don't need cache.
type noCacheConfig struct{}

func (n noCacheConfig) GetCacheTTL() time.Duration { return 0 }

// NoCacheConfig can be used when a process doesn't need to handle cache.
var NoCacheConfig CacheConfig = noCacheConfig{}

type EmailConfig interface {
	SMTPAddress() string
	GetSMTPFrom() string
}

// noEmailConfig is a null implementation for processes that don't need email.
type noEmailConfig struct{}

func (n noEmailConfig) SMTPAddress() string { return "" }
func (n noEmailConfig) GetSMTPFrom() string  { return "" }

// NoEmailConfig can be used when a process doesn't need to send emails.
var NoEmailConfig EmailConfig = noEmailConfig{}

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

// noFrontendConfig is a null implementation for processes that don't need frontend URL.
type noFrontendConfig struct{}

func (n noFrontendConfig) GetFrontendURL() *url.URL { return nil }

// NoFrontendConfig can be used when a process doesn't need to know the frontend URL.
var NoFrontendConfig FrontendConfig = noFrontendConfig{}

// noQueueConfig is a null implementation for processes that don't need a queue.
type noQueueConfig struct{}

func (n noQueueConfig) SQSConfig() (string, string, string) { return "", "", "" }

// NoQueueConfig can be used when a process doesn't need to initialize a queue.
var NoQueueConfig QueueConfig = noQueueConfig{}
