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
