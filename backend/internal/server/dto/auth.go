package dto

// Auth DTOs
// LoginRequest represents the request body for login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for login.
type LoginResponse struct {
	Success bool `json:"success"`
}

// MessageResponse represents a simple message response.
type MessageResponse struct {
	Message string `json:"message"`
}
