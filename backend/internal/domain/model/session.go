package model

import "time"

// SessionRole provides SessionRole related functionality.
type SessionRole string

const (
	// SessionRoleAdmin provides SessionRoleAdmin related functionality.
	SessionRoleAdmin SessionRole = "admin"
	// SessionRoleBuyer provides SessionRoleBuyer related functionality.
	SessionRoleBuyer SessionRole = "buyer"
)

// Session provides Session related functionality.
type Session struct {
	ID        string
	UserID    int
	Role      SessionRole
	CreatedAt time.Time
}
