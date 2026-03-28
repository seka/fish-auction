package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// 必要な環境変数の最小セットをデフォルトで持っておく
	defaultEnv := map[string]string{
		"SERVER_ADDRESS": ":8080",
		"DB_HOST":        "localhost",
		"DB_PORT":        "5432",
		"DB_USER":        "postgres",
		"DB_PASSWORD":    "postgres",
		"DB_NAME":        "fish_auction",
	}

	tests := []struct {
		name        string
		env         map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name:    "Success with default FRONTEND_URL",
			env:     map[string]string{},
			wantErr: false,
		},
		{
			name: "Success with valid FRONTEND_URL",
			env: map[string]string{
				"FRONTEND_URL": "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "Error with missing scheme in FRONTEND_URL",
			env: map[string]string{
				"FRONTEND_URL": "localhost:3000",
			},
			wantErr:     true,
			errContains: "invalid or missing FRONTEND_URL",
		},
		{
			name: "Error with missing host in FRONTEND_URL",
			env: map[string]string{
				"FRONTEND_URL": "https://",
			},
			wantErr:     true,
			errContains: "invalid or missing FRONTEND_URL",
		},
		{
			name: "Success with empty FRONTEND_URL (falls back to default)",
			env: map[string]string{
				"FRONTEND_URL": "",
			},
			wantErr: false,
		},
		{
			name: "Error with invalid URL characters",
			env: map[string]string{
				"FRONTEND_URL": "://invalid",
			},
			wantErr:     true,
			errContains: "invalid or missing FRONTEND_URL",
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

			cfg, err := Load()
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				if val, ok := tt.env["FRONTEND_URL"]; ok && val != "" {
					assert.Equal(t, val, cfg.FrontendURL.String())
				} else {
					assert.Equal(t, "https://localhost", cfg.FrontendURL.String())
				}
			}
		})
	}
}
