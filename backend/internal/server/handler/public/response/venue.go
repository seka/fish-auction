package response

import "time"

// Venue represents a public view of a venue.
type Venue struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
