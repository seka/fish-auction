package dto

import "time"

// Item DTOs
type CreateItemRequest struct {
	AuctionID   int    `json:"auction_id"`
	FishermanID int    `json:"fisherman_id"`
	FishType    string `json:"fish_type"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
}

type ItemResponse struct {
	ID          int       `json:"id"`
	AuctionID   int       `json:"auction_id"`
	FishermanID int       `json:"fisherman_id"`
	FishType    string    `json:"fish_type"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	Status      string    `json:"status"`
	HighestBid  *int      `json:"highest_bid,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
