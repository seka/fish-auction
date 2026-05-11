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
		name           string
		env            map[string]string
		wantFrontend   string
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

func TestValidateAppServerConfig(t *testing.T) {
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
			err := ValidateAppServerConfig(cfg)
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
