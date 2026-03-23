package model

import (
	"time"
)

// Bid provides Bid related functionality.
type Bid struct {
	ID        int
	ItemID    int
	BuyerID   int
	Price     BidPrice
	CreatedAt time.Time
}
