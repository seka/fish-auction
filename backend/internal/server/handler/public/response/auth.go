package response

// Login represents the login response.
type Login struct {
	Success bool `json:"success"`
}

// Buyer represents the buyer's basic session info for public auth endpoints.
type Buyer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
