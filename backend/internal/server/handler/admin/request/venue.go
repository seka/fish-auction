package request

// CreateVenue holds data for venue creation.
type CreateVenue struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

// UpdateVenue holds data for updating a venue.
type UpdateVenue struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}
