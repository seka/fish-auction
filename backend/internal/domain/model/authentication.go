package model

import "time"

type Authentication struct {
	ID             int
	BuyerID        int
	Email          string
	PasswordHash   string
	AuthType       string
	FailedAttempts int
	LockedUntil    *time.Time
	LastLoginAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
