package dto

import "time"

// CreateVenueRequest is a data transfer object.
type CreateVenueRequest struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

// UpdateVenueRequest represents the request body for updating a venue.
type UpdateVenueRequest struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

// VenueResponse represents the response body for a venue.
type VenueResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
