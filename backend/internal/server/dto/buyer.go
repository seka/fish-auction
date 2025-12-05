package dto

// Buyer DTOs
type CreateBuyerRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
	ContactInfo  string `json:"contact_info"`
}

type LoginBuyerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BuyerResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PurchaseResponse struct {
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
