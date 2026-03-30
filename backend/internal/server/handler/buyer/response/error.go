package response

// Error represents an error response for buyers.
type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
