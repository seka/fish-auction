package middleware

type contextKey string

const (
	AdminIDKey contextKey = "admin_id"
	BuyerIDKey contextKey = "buyer_id"
)
