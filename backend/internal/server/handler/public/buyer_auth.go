package public

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

// BuyerAuthHandler handles public HTTP requests related to buyer authentication.
type BuyerAuthHandler struct {
	loginUseCase buyer.LoginBuyerUseCase
	sessionRepo  repository.SessionRepository
}

// NewBuyerAuthHandler creates a new BuyerAuthHandler instance.
func NewBuyerAuthHandler(r registry.UseCase, sessionRepo repository.SessionRepository) *BuyerAuthHandler {
	return &BuyerAuthHandler{
		loginUseCase: r.NewLoginBuyerUseCase(),
		sessionRepo:  sessionRepo,
	}
}

// Login handles the buyer login request.
func (h *BuyerAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.Login
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	buy, err := h.loginUseCase.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	if buy == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID, err := h.sessionRepo.Create(r.Context(), buy.ID, model.SessionRoleBuyer)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	setSessionCookie(w, "buyer_session", sessionID)

	resp := response.Buyer{ID: buy.ID, Name: buy.Name}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Logout handles the buyer logout request.
func (h *BuyerAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("buyer_session"); err == nil {
		if err := h.sessionRepo.Delete(r.Context(), cookie.Value); err != nil {
			util.HandleError(w, err)
			return
		}
	}

	clearSessionCookie(w, "buyer_session")

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response.Message{Message: "Logged out"})
}

// RegisterRoutes registers the public buyer handler routes to the given mux.
func (h *BuyerAuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/buyer/login", h.Login)
	mux.HandleFunc("POST /api/buyer/logout", h.Logout)
}
