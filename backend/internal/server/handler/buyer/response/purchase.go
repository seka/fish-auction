package response

// Purchase represents a view of a purchase for the buyer.
type Purchase struct {
	ID          int    `json:"id"`
	ItemID      int    `json:"item_id"`
	FishType    string `json:"fish_type"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
	Price       int    `json:"price"`
	BuyerID     int    `json:"buyer_id"`
	AuctionID   int    `json:"auction_id"`
	AuctionDate string `json:"auction_date"`
	CreatedAt   string `json:"created_at"`
}
