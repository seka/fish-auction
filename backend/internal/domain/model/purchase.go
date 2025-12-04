package model

import "time"

// Purchase represents a buyer's purchase history item
type Purchase struct {
	ID          int
	ItemID      int
	FishType    string
	Quantity    int
	Unit        string
	Price       int
	BuyerID     int
	AuctionID   int
	AuctionDate string
	CreatedAt   time.Time
}
