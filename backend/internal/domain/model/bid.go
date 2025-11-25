package model

import "time"

type Bid struct {
	ID        int
	ItemID    int
	BuyerID   int
	Price     int
	CreatedAt time.Time
}
