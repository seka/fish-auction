package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSeedConfig(t *testing.T) {
	defaultEnv := map[string]string{
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_DB":       "fish_auction",
	}

	tests := []struct {
		name       string
		env        map[string]string
		wantAppEnv string
	}{
		{
			name:       "Default APP_ENV",
			env:        map[string]string{},
			wantAppEnv: "development",
		},
		{
			name: "Explicit APP_ENV",
			env: map[string]string{
				"APP_ENV": "test",
			},
			wantAppEnv: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range defaultEnv {
				t.Setenv(k, v)
			}
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg := NewSeedConfig()
			assert.NotNil(t, cfg)
			assert.Equal(t, "localhost", cfg.PostgresHost)
			assert.Equal(t, tt.wantAppEnv, cfg.AppEnv)
		})
	}
}

func TestSeedConfig_DBConnectionURL(t *testing.T) {
	cfg := &SeedConfig{
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
