package model

import "time"

// Venue represents a physical auction venue (会場)
type Venue struct {
	ID          int
	Name        string
	Location    string
	Description string
	CreatedAt   time.Time
}
