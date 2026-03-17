package dto

// UpdatePasswordRequest represents the request body for updating a password.
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
