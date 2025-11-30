package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type BuyerAuthMiddleware struct{}

func NewBuyerAuthMiddleware() *BuyerAuthMiddleware {
	return &BuyerAuthMiddleware{}
}

func (m *BuyerAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check session cookie
		cookie, err := r.Cookie("buyer_session")
		if err != nil || cookie.Value != "authenticated" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get buyer ID from cookie
		idCookie, err := r.Cookie("buyer_id")
		if err != nil {
			http.Error(w, "Unauthorized: missing buyer ID", http.StatusUnauthorized)
			return
		}

		buyerID, err := strconv.Atoi(idCookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized: invalid buyer ID", http.StatusUnauthorized)
			return
		}

		// Set buyer ID in context
		ctx := context.WithValue(r.Context(), "buyer_id", buyerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
