package handler

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
)

// AuthHandler handles HTTP requests related to authentication.
type AuthHandler struct {
	loginUseCase auth.LoginUseCase
	sessionRepo  repository.SessionRepository
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(r registry.UseCase, sessionRepo repository.SessionRepository) *AuthHandler {
	return &AuthHandler{
		loginUseCase: r.NewLoginUseCase(),
		sessionRepo:  sessionRepo,
	}
}

// Login handles the login request.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	admin, err := h.loginUseCase.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	if admin == nil {
		util.HandleError(w, &domainerrors.UnauthorizedError{Message: "Invalid credentials"})
		return
	}

	sessionID, err := h.sessionRepo.Create(r.Context(), admin.ID, model.SessionRoleAdmin)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	setSessionCookie(w, "admin_session", sessionID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"}); err != nil {
		util.HandleError(w, err)
		return
	}
}

// Logout handles the logout request.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("admin_session"); err == nil {
		if err := h.sessionRepo.Delete(r.Context(), cookie.Value); err != nil {
			util.HandleError(w, err)
			return
		}
	}

	clearSessionCookie(w, "admin_session")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"}); err != nil {
		util.HandleError(w, err)
		return
	}
}

// RegisterRoutes registers the auth handler routes to the given mux.
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
