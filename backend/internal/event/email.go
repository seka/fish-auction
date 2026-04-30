package event

// EmailType represents the type of email to be sent.
type EmailType string

const (
	EmailTypeBuyerPasswordReset EmailType = "buyer_password_reset"
	EmailTypeAdminPasswordReset EmailType = "admin_password_reset"
)

// EmailMessage is the wire format for email job messages.
type EmailMessage struct {
	EmailType EmailType `json:"email_type"`
	To        string    `json:"to"`
	ResetURL  string    `json:"reset_url,omitempty"`
}
