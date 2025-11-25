package dto

import "time"

// Item DTOs
type CreateItemRequest struct {
	FishermanID int    `json:"fisherman_id"`
	FishType    string `json:"fish_type"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
}

type ItemResponse struct {
	ID          int       `json:"id"`
	FishermanID int       `json:"fisherman_id"`
	FishType    string    `json:"fish_type"`
	Quantity    int       `json:"quantity"`
	Unit        string    `json:"unit"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
