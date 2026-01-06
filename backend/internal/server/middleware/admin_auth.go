package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AdminAuthMiddleware struct{}

func NewAdminAuthMiddleware() *AdminAuthMiddleware {
	return &AdminAuthMiddleware{}
}

func (m *AdminAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check session cookie
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Value != "authenticated" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized", "message": "Admin session cookie missing or invalid"})
			return
		}

		// Get admin ID from cookie
		idCookie, err := r.Cookie("admin_id")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized", "message": "Admin ID cookie missing"})
			return
		}

		var adminID int
		if _, err := fmt.Sscanf(idCookie.Value, "%d", &adminID); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized", "message": "Invalid Admin ID in cookie"})
			return
		}

		// Set admin ID in context
		ctx := context.WithValue(r.Context(), "admin_id", adminID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
