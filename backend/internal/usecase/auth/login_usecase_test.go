package auth_test

import (
	"context"
	"testing"

	"github.com/seka/fish-auction/backend/internal/usecase/auth"
)

func TestLoginUseCase_Execute(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "AdminPassword",
			password: "admin-password",
			want:     true,
		},
		{
			name:     "WrongPassword",
			password: "guest",
			want:     false,
		},
	}

	uc := auth.NewLoginUseCase()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.Execute(context.Background(), tt.password)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
