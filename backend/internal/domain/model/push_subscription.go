package model

import "time"

// PushSubscription represents a Web Push subscription
type PushSubscription struct {
	ID        int       `json:"id"`
	BuyerID   int       `json:"buyer_id"`
	Endpoint  string    `json:"endpoint"`
	P256dh    string    `json:"p256dh"`
	Auth      string    `json:"auth"`
	CreatedAt time.Time `json:"created_at"`
}
