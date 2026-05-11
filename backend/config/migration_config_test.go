package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMigrationConfig(t *testing.T) {
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

		cfg := NewMigrationConfig()
		assert.NotNil(t, cfg)
		assert.Equal(t, "localhost", cfg.PostgresHost)
		assert.Equal(t, "fish_auction", cfg.PostgresDB)
		assert.Equal(t, "disable", cfg.PostgresSslMode)
	})
}

func TestMigrationConfig_DBConnectionURL(t *testing.T) {
	cfg := &MigrationConfig{
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
