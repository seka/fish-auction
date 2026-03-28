package dto

// ResetPasswordRequest provides ResetPasswordRequest related functionality.
type ResetPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordVerifyRequest provides ResetPasswordVerifyRequest related functionality.
type ResetPasswordVerifyRequest struct {
	Token string `json:"token"`
}

// ResetPasswordConfirmRequest provides ResetPasswordConfirmRequest related functionality.
type ResetPasswordConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
