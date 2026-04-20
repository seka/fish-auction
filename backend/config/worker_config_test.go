package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadWorkerConfig(t *testing.T) {
	// 必要な環境変数の最小セット
	defaultEnv := map[string]string{
		"POSTGRES_HOST":     "localhost",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_DB":       "fish_auction",
		"VAPID_PUBLIC_KEY":  "test-public-key",
		"VAPID_PRIVATE_KEY": "test-private-key",
	}

	tests := []struct {
		name        string
		env         map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name:    "Success",
			env:     map[string]string{},
			wantErr: false,
		},
		{
			name: "Missing env doesn't cause load error",
			env: map[string]string{
				"POSTGRES_HOST": "",
			},
			wantErr: false,
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

			cfg, err := LoadWorkerConfig()
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, "localhost", cfg.RedisHost)
			}
		})
	}
}
