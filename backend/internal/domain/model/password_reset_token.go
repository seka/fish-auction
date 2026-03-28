package model

import "time"

// PasswordResetToken represents a password reset token entry in the system.
type PasswordResetToken struct {
	UserID    int
	Role      string
	TokenHash string
	ExpiresAt time.Time
}
