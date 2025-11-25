package dto

// Buyer DTOs
type CreateBuyerRequest struct {
	Name string `json:"name"`
}

type BuyerResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
