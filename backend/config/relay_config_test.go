package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRelayConfig(t *testing.T) {
	defaultEnv := map[string]string{
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_DB":       "fish_auction",
	}

	t.Run("Success", func(t *testing.T) {
		os.Clearenv()
		for k, v := range defaultEnv {
			t.Setenv(k, v)
		}

		cfg := NewRelayConfig()
		assert.NotNil(t, cfg)
		assert.Equal(t, "localhost", cfg.PostgresHost)
		assert.Equal(t, "ap-northeast-1", cfg.SQSRegion)
		assert.Equal(t, "http://localhost:4566/000000000000/notification-queue", cfg.SQSQueueURL)
	})
}

func TestRelayConfig_DBConnectionURL(t *testing.T) {
	cfg := &RelayConfig{
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

func TestRelayConfig_SQSConfig(t *testing.T) {
	cfg := &RelayConfig{
		SQSRegion:   "ap-northeast-1",
		SQSQueueURL: "https://sqs.example.com/queue",
		SQSEndpoint: "https://sqs.example.com",
	}

	region, queueURL, endpoint := cfg.SQSConfig()
	assert.Equal(t, "ap-northeast-1", region)
	assert.Equal(t, "https://sqs.example.com/queue", queueURL)
	assert.Equal(t, "https://sqs.example.com", endpoint)
}
