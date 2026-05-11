package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInitAdminConfig(t *testing.T) {
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

		cfg := NewInitAdminConfig()
		assert.NotNil(t, cfg)
		assert.Equal(t, "localhost", cfg.PostgresHost)
		assert.Equal(t, "fish_auction", cfg.PostgresDB)
		assert.Equal(t, "disable", cfg.PostgresSslMode)
	})
}

func TestInitAdminConfig_DBConnectionURL(t *testing.T) {
	cfg := &InitAdminConfig{
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

func TestInitAdminConfig_Validate(t *testing.T) {
	full := &InitAdminConfig{
		PostgresHost:     "localhost",
		PostgresPort:     "5432",
		PostgresUser:     "postgres",
		PostgresPassword: "postgres",
		PostgresDB:       "fish_auction",
		PostgresSslMode:  "disable",
	}

	tests := []struct {
		name        string
		mutate      func(*InitAdminConfig)
		wantErr     bool
		errContains string
	}{
		{
			name:    "All required env present",
			mutate:  func(c *InitAdminConfig) {},
			wantErr: false,
		},
		{
			name:        "Missing POSTGRES_HOST",
			mutate:      func(c *InitAdminConfig) { c.PostgresHost = "" },
			wantErr:     true,
			errContains: "POSTGRES_HOST",
		},
		{
			name: "Multiple missing values are all reported",
			mutate: func(c *InitAdminConfig) {
				c.PostgresUser = ""
				c.PostgresDB = ""
			},
			wantErr:     true,
			errContains: "POSTGRES_USER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := *full
			tt.mutate(&cfg)
			err := cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
