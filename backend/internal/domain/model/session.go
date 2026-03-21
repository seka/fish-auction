package model

import "time"

type SessionRole string

const (
	SessionRoleAdmin SessionRole = "admin"
	SessionRoleBuyer SessionRole = "buyer"
)

type Session struct {
	ID        string      `json:"id"`
	UserID    int         `json:"user_id"`
	Role      SessionRole `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
}
