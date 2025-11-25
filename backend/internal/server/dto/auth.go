package dto

// Auth DTOs
type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool `json:"success"`
}
