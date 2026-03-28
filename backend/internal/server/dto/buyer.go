package dto

// CreateBuyerRequest is a data transfer object.
type CreateBuyerRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
	ContactInfo  string `json:"contact_info"`
}

// LoginBuyerRequest represents the request body for buyer login.
type LoginBuyerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// BuyerResponse represents the response body for a buyer.
type BuyerResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PurchaseResponse represents the response body for a purchase.
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

// BuyerMeResponse represents the response body for the current authenticated buyer.
type BuyerMeResponse struct {
	Authenticated bool   `json:"authenticated"`
	BuyerID       int    `json:"buyer_id"`
	Name          string `json:"name"`
}
