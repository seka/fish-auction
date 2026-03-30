package response

// Fisherman represents a view of a fisherman for admins.
type Fisherman struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
