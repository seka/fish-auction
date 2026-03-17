package dto

import "time"

// Venue DTOs
// CreateVenueRequest represents the request body for creating a venue.
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
