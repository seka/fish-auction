package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppServerConfig(t *testing.T) {
	defaultEnv := map[string]string{
		"SERVER_PORT":       "8080",
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_DB":       "fish_auction",
	}

	tests := []struct {
		name         string
		env          map[string]string
		wantFrontend string
	}{
		{
			name:         "Default FRONTEND_URL",
			env:          map[string]string{},
			wantFrontend: "https://localhost",
		},
		{
			name: "Valid FRONTEND_URL",
			env: map[string]string{
				"FRONTEND_URL": "https://example.com",
			},
			wantFrontend: "https://example.com",
		},
		{
			name: "Empty FRONTEND_URL falls back to default",
			env: map[string]string{
				"FRONTEND_URL": "",
			},
			wantFrontend: "https://localhost",
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

			cfg := NewAppServerConfig()
			assert.NotNil(t, cfg)
			assert.Equal(t, tt.wantFrontend, cfg.GetFrontendURL().String())
		})
	}
}

func TestAppServerConfig_Validate(t *testing.T) {
	defaultEnv := map[string]string{
		"SERVER_PORT":       "8080",
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_DB":       "fish_auction",
	}

	tests := []struct {
		name        string
		env         map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name:    "Valid default",
			env:     map[string]string{},
			wantErr: false,
		},
		{
			name: "Valid FRONTEND_URL",
			env: map[string]string{
				"FRONTEND_URL": "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "Missing scheme in FRONTEND_URL",
			env: map[string]string{
				"FRONTEND_URL": "localhost:3000",
			},
			wantErr:     true,
			errContains: "invalid FRONTEND_URL",
		},
		{
			name: "Missing host in FRONTEND_URL",
			env: map[string]string{
				"FRONTEND_URL": "https://",
			},
			wantErr:     true,
			errContains: "invalid FRONTEND_URL",
		},
		{
			name: "Invalid URL characters",
			env: map[string]string{
				"FRONTEND_URL": "://invalid",
			},
			wantErr:     true,
			errContains: "invalid FRONTEND_URL",
		},
		{
			name: "Valid TRUSTED_PROXIES",
			env: map[string]string{
				"TRUSTED_PROXIES": "10.0.0.0/16, 192.168.0.0/24",
			},
			wantErr: false,
		},
		{
			name: "Empty TRUSTED_PROXIES is allowed",
			env: map[string]string{
				"TRUSTED_PROXIES": "",
			},
			wantErr: false,
		},
		{
			name: "Invalid TRUSTED_PROXIES CIDR",
			env: map[string]string{
				"TRUSTED_PROXIES": "10.0.0.0/16,not-a-cidr",
			},
			wantErr:     true,
			errContains: "invalid TRUSTED_PROXIES",
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

			cfg := NewAppServerConfig()
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

func TestAppServerConfig_ServerAddr(t *testing.T) {
	cfg := &AppServerConfig{ServerHost: "0.0.0.0", ServerPort: "8080"}
	assert.Equal(t, "0.0.0.0:8080", cfg.ServerAddr())
}

func TestAppServerConfig_RedisAddr(t *testing.T) {
	cfg := &AppServerConfig{RedisHost: "redis.example.com", RedisPort: "6379"}
	assert.Equal(t, "redis.example.com:6379", cfg.RedisAddr())
}

func TestAppServerConfig_SMTPAddress(t *testing.T) {
	cfg := &AppServerConfig{SMTPHost: "smtp.example.com", SMTPPort: "1025"}
	assert.Equal(t, "smtp.example.com:1025", cfg.SMTPAddress())
}

func TestAppServerConfig_DBConnectionURL(t *testing.T) {
	cfg := &AppServerConfig{
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
