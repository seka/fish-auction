package middleware

import (
	"log/slog"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/util"
)

// AdminAuthMiddleware provides AdminAuthMiddleware related functionality.
type AdminAuthMiddleware struct {
	sessionRepo repository.SessionRepository
}

// NewAdminAuthMiddleware creates a new AdminAuthMiddleware instance.
func NewAdminAuthMiddleware(sessionRepo repository.SessionRepository) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{sessionRepo: sessionRepo}
}

// Handle provides Handle related functionality.
func (m *AdminAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Value == "" {
			slog.Warn("auth: admin cookie missing",
				"remote_addr", r.RemoteAddr,
				"path", r.URL.Path,
				"request_id", RequestIDFromContext(r.Context()),
			)
			util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		session, err := m.sessionRepo.FindByID(r.Context(), cookie.Value)
		if err != nil {
			slog.Error("auth: admin session lookup failed",
				"err", err,
				"request_id", RequestIDFromContext(r.Context()),
			)
			util.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if session == nil || session.Role != model.SessionRoleAdmin {
			slog.Warn("auth: admin session invalid",
				"remote_addr", r.RemoteAddr,
				"has_session", session != nil,
				"request_id", RequestIDFromContext(r.Context()),
			)
			util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := WithAdminID(r.Context(), session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
