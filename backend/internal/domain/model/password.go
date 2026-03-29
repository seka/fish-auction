package model

import (
	"errors"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
)

// Password represents a raw password that meets complexity requirements.
// It is primarily used when setting or updating a password.
type Password struct {
	value string
}

// NewPassword creates a new Password after validating its complexity.
// Complexity rules:
// - At least 8 characters long
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one digit
func NewPassword(v string) (*Password, error) {
	if err := validateComplexity(v); err != nil {
		return nil, err
	}
	return &Password{value: v}, nil
}

// Hash generates a bcrypt hash of the password.
func (p *Password) Hash() (HashedPassword, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return HashedPassword{}, fmt.Errorf("failed to hash password: %w", err)
	}
	return NewHashedPassword(string(hashedBytes)), nil
}

// String returns a masked representation of the password.
func (p *Password) String() string {
	return "********"
}

// HashedPassword represents a password hash stored in the database.
// It is used for verifying incoming raw passwords against the stored hash.
type HashedPassword struct {
	value string
}

// NewHashedPassword creates a HashedPassword from an existing hash string.
func NewHashedPassword(h string) HashedPassword {
	return HashedPassword{value: h}
}

// Verify compares a raw password with the hashed password.
func (hp HashedPassword) Verify(raw string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hp.value), []byte(raw)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return &domainErrors.UnauthorizedError{Message: "Invalid credentials"}
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}
	return nil
}

// String returns a label indicating this is a hashed value.
func (hp HashedPassword) String() string {
	return "[hashed]"
}

// Raw returns the underlying hash string.
func (hp HashedPassword) Raw() string {
	return hp.value
}

func validateComplexity(p string) error {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		minLen     = 8
		maxLen     = 72 // bcrypt limit
	)

	if len(p) < minLen || len(p) > maxLen {
		return &domainErrors.ValidationError{
			Field:   "password",
			Message: fmt.Sprintf("password must be between %d and %d characters long", minLen, maxLen),
		}
	}

	for _, char := range p {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return &domainErrors.ValidationError{
			Field:   "password",
			Message: "password must contain at least one uppercase letter, one lowercase letter, and one number",
		}
	}

	return nil
}
