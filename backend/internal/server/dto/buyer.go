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
