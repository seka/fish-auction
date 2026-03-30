package middleware

import (
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/util"
)

// BuyerAuthMiddleware provides BuyerAuthMiddleware related functionality.
type BuyerAuthMiddleware struct {
	sessionRepo repository.SessionRepository
}

// NewBuyerAuthMiddleware creates a new BuyerAuthMiddleware instance.
func NewBuyerAuthMiddleware(sessionRepo repository.SessionRepository) *BuyerAuthMiddleware {
	return &BuyerAuthMiddleware{sessionRepo: sessionRepo}
}

// Handle provides Handle related functionality.
func (m *BuyerAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("buyer_session")
		if err != nil || cookie.Value == "" {
			util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		session, err := m.sessionRepo.FindByID(r.Context(), cookie.Value)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if session == nil || session.Role != model.SessionRoleBuyer {
			util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := WithBuyerID(r.Context(), session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
