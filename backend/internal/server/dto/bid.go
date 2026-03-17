package dto

import "time"

// Bid/Transaction DTOs
// CreateBidRequest represents the request body for creating a bid.
type CreateBidRequest struct {
	ItemID  int `json:"item_id"`
	BuyerID int `json:"buyer_id"`
	Price   int `json:"price"`
}

// BidResponse represents the response body for a bid.
type BidResponse struct {
	ID        int       `json:"id"`
	ItemID    int       `json:"item_id"`
	BuyerID   int       `json:"buyer_id"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
