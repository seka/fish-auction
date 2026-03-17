package dto

// Fisherman DTOs
// CreateFishermanRequest represents the request body for creating a fisherman.
type CreateFishermanRequest struct {
	Name string `json:"name"`
}

// FishermanResponse represents the response body for a fisherman.
type FishermanResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
