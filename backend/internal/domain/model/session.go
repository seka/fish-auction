package model

import "time"

type SessionRole string

const (
	SessionRoleAdmin SessionRole = "admin"
	SessionRoleBuyer SessionRole = "buyer"
)

type Session struct {
	ID        string
	UserID    int
	Role      SessionRole
	CreatedAt time.Time
}
