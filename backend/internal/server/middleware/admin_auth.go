package middleware

import (
	"context"
	"fmt"
	"net/http"
)

type AdminAuthMiddleware struct{}

func NewAdminAuthMiddleware() *AdminAuthMiddleware {
	return &AdminAuthMiddleware{}
}

func (m *AdminAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check session cookie
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Value != "authenticated" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get admin ID from cookie
		idCookie, err := r.Cookie("admin_id")
		if err != nil {
			http.Error(w, "Unauthorized: Admin ID missing", http.StatusUnauthorized)
			return
		}

		var adminID int
		if _, err := fmt.Sscanf(idCookie.Value, "%d", &adminID); err != nil {
			http.Error(w, "Unauthorized: invalid Admin ID", http.StatusUnauthorized)
			return
		}

		// Set admin ID in context
		ctx := context.WithValue(r.Context(), "admin_id", adminID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
