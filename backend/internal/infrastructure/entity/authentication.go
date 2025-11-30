package entity

import (
	"strings"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type Authentication struct {
	ID             int        `db:"id"`
	BuyerID        int        `db:"buyer_id"`
	Email          string     `db:"email"`
	PasswordHash   string     `db:"password_hash"`
	AuthType       string     `db:"auth_type"`
	FailedAttempts int        `db:"failed_attempts"`
	LockedUntil    *time.Time `db:"locked_until"`
	LastLoginAt    *time.Time `db:"last_login_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

func (a *Authentication) Validate() error {
	if strings.TrimSpace(a.Email) == "" {
		return &apperrors.ValidationError{Field: "email", Message: "Email is required"}
	}
	if a.PasswordHash == "" {
		return &apperrors.ValidationError{Field: "password_hash", Message: "Password hash is required"}
	}
	if a.BuyerID <= 0 {
		return &apperrors.ValidationError{Field: "buyer_id", Message: "Buyer ID must be positive"}
	}
	return nil
}

func (a *Authentication) IsLocked() bool {
	if a.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*a.LockedUntil)
}

func (a *Authentication) ToModel() *model.Authentication {
	return &model.Authentication{
		ID:             a.ID,
		BuyerID:        a.BuyerID,
		Email:          a.Email,
		PasswordHash:   a.PasswordHash,
		AuthType:       a.AuthType,
		FailedAttempts: a.FailedAttempts,
		LockedUntil:    a.LockedUntil,
		LastLoginAt:    a.LastLoginAt,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}
