package model

import "time"

// Admin represents an administrative user
type Admin struct {
	ID             int
	Email          string
	PasswordHash   string
	FailedAttempts int
	LockedUntil    *time.Time
	CreatedAt      time.Time
}
