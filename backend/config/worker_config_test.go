package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkerConfig(t *testing.T) {
	defaultEnv := map[string]string{
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_DB":       "fish_auction",
		"VAPID_PUBLIC_KEY":  "test-public-key",
		"VAPID_PRIVATE_KEY": "test-private-key",
	}

	t.Run("Success", func(t *testing.T) {
		os.Clearenv()
		for k, v := range defaultEnv {
			t.Setenv(k, v)
		}

		cfg := NewWorkerConfig()
		assert.NotNil(t, cfg)
		assert.Equal(t, "localhost", cfg.RedisHost)
	})
}

func TestWorkerConfig_RedisAddr(t *testing.T) {
	cfg := &WorkerConfig{RedisHost: "redis.example.com", RedisPort: "6379"}
	assert.Equal(t, "redis.example.com:6379", cfg.RedisAddr())
}

func TestWorkerConfig_SMTPAddress(t *testing.T) {
	cfg := &WorkerConfig{SMTPHost: "smtp.example.com", SMTPPort: "1025"}
	assert.Equal(t, "smtp.example.com:1025", cfg.SMTPAddress())
}

func TestWorkerConfig_VAPIDConfig(t *testing.T) {
	cfg := &WorkerConfig{
		VAPIDPublicKey:  "public",
		VAPIDPrivateKey: "private",
		VAPIDSubject:    "mailto:admin@example.com",
	}

	pub, priv, sub := cfg.VAPIDConfig()
	assert.Equal(t, "public", pub)
	assert.Equal(t, "private", priv)
	assert.Equal(t, "mailto:admin@example.com", sub)
}

func TestWorkerConfig_SQSConfig(t *testing.T) {
	cfg := &WorkerConfig{
		SQSRegion:   "ap-northeast-1",
		SQSQueueURL: "https://sqs.example.com/queue",
		SQSEndpoint: "https://sqs.example.com",
	}

	region, queueURL, endpoint := cfg.SQSConfig()
	assert.Equal(t, "ap-northeast-1", region)
	assert.Equal(t, "https://sqs.example.com/queue", queueURL)
	assert.Equal(t, "https://sqs.example.com", endpoint)
}

func TestWorkerConfig_DBConnectionURL(t *testing.T) {
	cfg := &WorkerConfig{
		PostgresHost:     "db.example.com",
		PostgresPort:     "5432",
		PostgresUser:     "user",
		PostgresPassword: "pass",
		PostgresDB:       "fish_auction",
		PostgresSslMode:  "require",
	}

	got := cfg.DBConnectionURL()
	want := "host=db.example.com port=5432 user=user password=pass dbname=fish_auction sslmode=require"
	assert.Equal(t, want, got)
}
