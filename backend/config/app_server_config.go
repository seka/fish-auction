package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// AppServerConfig represents the configuration for the API server.
type AppServerConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSslMode  string
	ServerHost       string
	ServerPort       string
	RedisHost        string
	RedisPort        string
	RedisDB          int
	CacheTTL         time.Duration
	SessionTTL       time.Duration
	AppEnv           string
	AllowedOrigins   string
	TrustedProxies   string
	SMTPHost         string
	SMTPPort         string
	SMTPFrom         string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	FrontendURL      *url.URL
}

// NewAppServerConfig は API サーバ用の設定を環境変数からロードする。
//
// 本関数は値の妥当性を検証しない。FRONTEND_URL のパースに失敗した場合は
// FrontendURL を nil のまま返すため、呼び出し側は必ず (*AppServerConfig).Validate
// を呼んで検証すること。サーバ起動経路（cmd/server, cmd/testing）はこれを遵守する。
func NewAppServerConfig() *AppServerConfig {
	cacheTTL := GetEnvInt("CACHE_TTL_SECONDS", 300)
	sessionTTL := GetEnvInt("SESSION_TTL_SECONDS", 86400)

	frontendURL, _ := url.Parse(GetEnv("FRONTEND_URL", "https://localhost"))

	return &AppServerConfig{
		PostgresHost:     GetEnv("POSTGRES_HOST", ""),
		PostgresPort:     GetEnv("POSTGRES_PORT", ""),
		PostgresUser:     GetEnv("POSTGRES_USER", ""),
		PostgresPassword: GetEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       GetEnv("POSTGRES_DB", ""),
		PostgresSslMode:  GetEnv("POSTGRES_SSLMODE", "disable"),
		ServerHost:       GetEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:       GetEnv("SERVER_PORT", "8080"),
		RedisHost:        GetEnv("REDIS_HOST", "localhost"),
		RedisPort:        GetEnv("REDIS_PORT", "6379"),
		RedisDB:          GetEnvInt("REDIS_DB", 0),
		CacheTTL:         time.Duration(cacheTTL) * time.Second,
		SessionTTL:       time.Duration(sessionTTL) * time.Second,
		AppEnv:           GetEnv("APP_ENV", "develop"),
		AllowedOrigins:   GetEnv("ALLOWED_ORIGINS", "https://localhost,http://localhost:3000"),
		TrustedProxies:   GetEnv("TRUSTED_PROXIES", ""),
		SMTPHost:         GetEnv("SMTP_HOST", "mailhog"),
		SMTPPort:         GetEnv("SMTP_PORT", "1025"),
		SMTPFrom:         GetEnv("SMTP_FROM", "noreply@fish-auction.com"),
		ReadTimeout:      time.Duration(GetEnvInt("SERVER_READ_TIMEOUT_SEC", 60)) * time.Second,
		WriteTimeout:     time.Duration(GetEnvInt("SERVER_WRITE_TIMEOUT_SEC", 60)) * time.Second,
		IdleTimeout:      time.Duration(GetEnvInt("SERVER_IDLE_TIMEOUT_SEC", 60)) * time.Second,
		FrontendURL:      frontendURL,
	}
}

// Validate は値の妥当性を検証する。
func (c *AppServerConfig) Validate() error {
	if c.FrontendURL == nil || c.FrontendURL.Scheme == "" || c.FrontendURL.Host == "" {
		return errors.New("invalid FRONTEND_URL: missing scheme or host")
	}
	for _, raw := range strings.Split(c.TrustedProxies, ",") {
		cidr := strings.TrimSpace(raw)
		if cidr == "" {
			continue
		}
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			return fmt.Errorf("invalid TRUSTED_PROXIES: %q is not a valid CIDR: %w", cidr, err)
		}
	}
	if err := validateSSLMode(c.AppEnv, c.PostgresSslMode); err != nil {
		return err
	}
	return nil
}

func (c *AppServerConfig) ServerAddr() string {
	return net.JoinHostPort(c.ServerHost, c.ServerPort)
}

func (c *AppServerConfig) RedisAddr() string {
	return net.JoinHostPort(c.RedisHost, c.RedisPort)
}

func (c *AppServerConfig) GetRedisDB() int {
	return c.RedisDB
}

func (c *AppServerConfig) GetSessionTTL() time.Duration {
	return c.SessionTTL
}

func (c *AppServerConfig) GetCacheTTL() time.Duration {
	return c.CacheTTL
}

func (c *AppServerConfig) SMTPAddress() string {
	return fmt.Sprintf("%s:%s", c.SMTPHost, c.SMTPPort)
}

func (c *AppServerConfig) GetSMTPFrom() string {
	return c.SMTPFrom
}

func (c *AppServerConfig) DBConnectionURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresSslMode)
}

func (c *AppServerConfig) GetFrontendURL() *url.URL {
	return c.FrontendURL
}
