package response

// Error represents an error response for admins.
type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
