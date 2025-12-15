package entity

import "time"

// Admin represents an administrative user
type Admin struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
