package middleware

import "context"

type contextKey string

const (
	// AdminIDKey provides AdminIDKey related functionality.
	AdminIDKey contextKey = "admin_id"
	// BuyerIDKey provides BuyerIDKey related functionality.
	BuyerIDKey contextKey = "buyer_id"
)

// AdminIDFromContext provides AdminIDFromContext related functionality.
func AdminIDFromContext(ctx context.Context) (int, bool) {
	adminID, ok := ctx.Value(AdminIDKey).(int)
	return adminID, ok
}

// BuyerIDFromContext provides BuyerIDFromContext related functionality.
func BuyerIDFromContext(ctx context.Context) (int, bool) {
	buyerID, ok := ctx.Value(BuyerIDKey).(int)
	return buyerID, ok
}
