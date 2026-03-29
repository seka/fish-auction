package model

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	// minPasswordLength is the minimum required length for a password.
	minPasswordLength = 8
	// maxPasswordLength is the maximum allowed length for a password.
	// 72 is chosen because bcrypt ignores any characters after the 72nd byte.
	maxPasswordLength = 72
)

// Password represents a raw password string (Value Object) that is guaranteed to be valid.
type Password string

// NewPassword validates and creates a new Password instance.
func NewPassword(v string) (Password, error) {
	p := Password(v)
	if err := p.Validate(); err != nil {
		return "", err
	}
	return p, nil
}

// Validate checks if the password satisfies complexity rules.
func (p Password) Validate() error {
	v := string(p)
	if len(v) < minPasswordLength {
		return &errors.ValidationError{
			Field:   "password",
			Message: fmt.Sprintf("must be at least %d characters long", minPasswordLength),
		}
	}
	if len(v) > maxPasswordLength {
		return &errors.ValidationError{
			Field:   "password",
			Message: fmt.Sprintf("must be no more than %d characters long", maxPasswordLength),
		}
	}

	hasUpper := strings.ContainsFunc(v, unicode.IsUpper)
	hasLower := strings.ContainsFunc(v, unicode.IsLower)
	hasDigit := strings.ContainsFunc(v, unicode.IsDigit)

	if !hasUpper || !hasLower || !hasDigit {
		return &errors.ValidationError{
			Field:   "password",
			Message: "must contain at least one uppercase letter, one lowercase letter, and one digit",
		}
	}

	return nil
}

// Hash returns the bcrypt hash of the password.
func (p Password) Hash() (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CompareWithHash verifies the password against a bcrypt hash.
func (p Password) CompareWithHash(hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(p))
}

// String implements fmt.Stringer to mask the password in logs.
func (p Password) String() string {
	return "********"
}
