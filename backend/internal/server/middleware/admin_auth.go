package middleware

import (
	"context"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type AdminAuthMiddleware struct {
	sessionRepo repository.SessionRepository
}

func NewAdminAuthMiddleware(sessionRepo repository.SessionRepository) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{sessionRepo: sessionRepo}
}

func (m *AdminAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session, err := m.sessionRepo.FindByID(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if session == nil || session.Role != model.SessionRoleAdmin {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AdminIDKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
