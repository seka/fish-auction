package request

// CreateBid holds data for placing a bid.
type CreateBid struct {
	ItemID int `json:"item_id"`
	Price  int `json:"price"`
}
