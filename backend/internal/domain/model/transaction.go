package model

import "time"

type Transaction struct {
	ID        int
	ItemID    int
	BuyerID   int
	Price     int
	CreatedAt time.Time
}
