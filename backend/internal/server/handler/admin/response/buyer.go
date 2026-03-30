package response

// Buyer represents a summary view of a buyer for admins.
type Buyer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
