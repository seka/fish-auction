package request

// UpdatePassword holds data for updating the admin's password.
type UpdatePassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
