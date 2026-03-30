package request

// ResetPassword holds data for requesting an admin password reset.
type ResetPassword struct {
	Email string `json:"email"`
}

// ResetPasswordVerify holds data for verifying an admin reset token.
type ResetPasswordVerify struct {
	Token string `json:"token"`
}

// ResetPasswordConfirm holds data for confirming an admin password reset.
type ResetPasswordConfirm struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
