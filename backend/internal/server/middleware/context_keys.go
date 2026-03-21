package middleware

import "context"

type contextKey string

const (
	AdminIDKey contextKey = "admin_id"
	BuyerIDKey contextKey = "buyer_id"
)

func AdminIDFromContext(ctx context.Context) (int, bool) {
	adminID, ok := ctx.Value(AdminIDKey).(int)
	return adminID, ok
}

func BuyerIDFromContext(ctx context.Context) (int, bool) {
	buyerID, ok := ctx.Value(BuyerIDKey).(int)
	return buyerID, ok
}
