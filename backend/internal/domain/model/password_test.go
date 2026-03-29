package model

import (
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		v       string
		wantErr bool
	}{
		{
			name:    "valid password",
			v:       "Password123",
			wantErr: false,
		},
		{
			name:    "too short",
			v:       "Pass12",
			wantErr: true,
		},
		{
			name:    "no uppercase",
			v:       "password123",
			wantErr: true,
		},
		{
			name:    "no lowercase",
			v:       "PASSWORD123",
			wantErr: true,
		},
		{
			name:    "no digit",
			v:       "Password",
			wantErr: true,
		},
		{
			name:    "too long",
			v:       "Password123Password123Password123Password123Password123Password123Password123", // 77 chars
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPassword(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if _, ok := err.(*errors.ValidationError); !ok {
					t.Errorf("NewPassword() error type = %T, want *errors.ValidationError", err)
				}
			}
		})
	}
}

func TestPassword_String(t *testing.T) {
	p, _ := NewPassword("Secret123")
	if p.String() != "********" {
		t.Errorf("Password.String() = %v, want ********", p.String())
	}
}

func TestPassword_Hash(t *testing.T) {
	p, _ := NewPassword("Password123")
	hash, err := p.Hash()
	if err != nil {
		t.Fatalf("Password.Hash() error = %v", err)
	}
	if len(hash) == 0 {
		t.Error("Password.Hash() returned empty hash")
	}
}
