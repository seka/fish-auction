package dto

import "time"

// Bid/Transaction DTOs
type CreateBidRequest struct {
	ItemID  int `json:"item_id"`
	BuyerID int `json:"buyer_id"`
	Price   int `json:"price"`
}

type BidResponse struct {
	ID        int       `json:"id"`
	ItemID    int       `json:"item_id"`
	BuyerID   int       `json:"buyer_id"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
