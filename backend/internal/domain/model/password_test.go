package model

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
)

func TestNewPassword_Complexity(t *testing.T) {
	tests := []struct {
		name    string
		v       string
		wantErr bool
	}{
		{"Too short", "Ab1", true},
		{"No uppercase", "password123", true},
		{"No lowercase", "PASSWORD123", true},
		{"No number", "Password", true},
		{"Valid complexity", "Admin123", false},
		{"Includes Japanese", "Admin123パス", true},
		{"Includes Emoji", "Admin123😀", true},
		{"Includes full-width", "Ａdmin123", true}, // Full-width 'A'
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPassword(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPassword_Masking(t *testing.T) {
	p, _ := NewPassword("Admin123")
	if p.String() != "********" {
		t.Errorf("Password.String() = %v, want %v", p.String(), "********")
	}
}

func TestHashedPassword_Verify(t *testing.T) {
	raw := "Admin123"
	p, _ := NewPassword(raw)
	hp, err := p.Hash()
	if err != nil {
		t.Fatalf("Failed to hash: %v", err)
	}

	// 1. Success
	if err := hp.Verify(raw); err != nil {
		t.Errorf("Verify() failed for valid password: %v", err)
	}

	// 2. Failure (wrong password)
	if err := hp.Verify("Wrong123"); err == nil {
		t.Error("Verify() should have failed for wrong password")
	} else {
		// Expect UnauthorizedError
		var unauthErr *domainErrors.UnauthorizedError
		if !errors.As(err, &unauthErr) {
			t.Errorf("Expected UnauthorizedError, got %v", err)
		}
	}

	// 3. Masking
	if hp.String() != "[hashed]" {
		t.Errorf("HashedPassword.String() = %v, want %v", hp.String(), "[hashed]")
	}
}

func TestHashedPassword_VerifyWithoutComplexityCheck(t *testing.T) {
	// 複雑性要件を満たさない既存パスワード（ポリシー導入前に登録されたユーザーを想定）を
	// bcrypt で直接ハッシュ化し、Verify が複雑性に関係なく照合できることを確認する
	simple := "simple"
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(simple), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}

	hp := NewHashedPassword(string(hashedBytes))
	if err := hp.Verify(simple); err != nil {
		t.Errorf("Verify() failed for pre-policy simple password: %v", err)
	}
}
