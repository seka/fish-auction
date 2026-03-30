package response

import "time"

// Bid represents a bid view for the buyer.
type Bid struct {
	ID        int       `json:"id"`
	ItemID    int       `json:"item_id"`
	BuyerID   int       `json:"buyer_id"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
