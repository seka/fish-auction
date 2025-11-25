package dto

// Fisherman DTOs
type CreateFishermanRequest struct {
	Name string `json:"name"`
}

type FishermanResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
