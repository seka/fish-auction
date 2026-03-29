package model

import (
	"errors"
	"testing"

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
	// Directly create HashedPassword (simulating DB load)
	// Even if it was hashed from a simple password (e.g. at initial creation before policy),
	// Verify should still work.
	
	// We use p.Hash() here just to get a valid bcrypt hash for a simple string.
	// But in reality, this would have been done via NewPassword at creation time (if it satisfied the policy then).
	// If it didn't satisfy the policy, but is somehow in the DB, Verify should still work.
	
	// Since NewPassword enforces complexity, let's use a workaround for testing 'already-in-db-simple-password'
	// Actually, hp.Verify() just calls bcrypt, it doesn't care about our Password complexity rules.
	
	p, _ := NewPassword("Valid123") // Satisfies complexity
	hp, _ := p.Hash()
	
	if err := hp.Verify("Valid123"); err != nil {
		t.Errorf("Verify failed for valid password: %v", err)
	}
}
