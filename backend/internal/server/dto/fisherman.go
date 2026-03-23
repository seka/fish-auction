package dto


// CreateFishermanRequest is a data transfer object.
type CreateFishermanRequest struct {
	Name string `json:"name"`
}

// FishermanResponse represents the response body for a fisherman.
type FishermanResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
