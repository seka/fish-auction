package request

// Login holds login request data.
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ResetPassword holds data for requesting a password reset.
type ResetPassword struct {
	Email string `json:"email"`
}

// ResetPasswordVerify holds data for verifying a reset token.
type ResetPasswordVerify struct {
	Token string `json:"token"`
}

// ResetPasswordConfirm holds data for confirming a password reset.
type ResetPasswordConfirm struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
