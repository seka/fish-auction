package public

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
)

// AdminAuthHandler handles public HTTP requests related to authentication.
type AdminAuthHandler struct {
	loginUseCase auth.LoginUseCase
	sessionRepo  repository.SessionRepository
}

// NewAdminAuthHandler creates a new AdminAuthHandler instance.
func NewAdminAuthHandler(r registry.UseCase, sessionRepo repository.SessionRepository) *AdminAuthHandler {
	return &AdminAuthHandler{
		loginUseCase: r.NewLoginUseCase(),
		sessionRepo:  sessionRepo,
	}
}

// Login handles the login request for admins.
func (h *AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.Login
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

	util.WriteJSON(w, http.StatusOK, response.Message{Message: "Login successful"})
}

// Logout handles the logout request for admins.
func (h *AdminAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("admin_session"); err == nil {
		if err := h.sessionRepo.Delete(r.Context(), cookie.Value); err != nil {
			util.HandleError(w, err)
			return
		}
	}

	clearSessionCookie(w, "admin_session")

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response.Message{Message: "Logged out"})
}

// RegisterRoutes registers the auth handler routes to the given mux.
func (h *AdminAuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/login", h.Login)
	mux.HandleFunc("POST /api/admin/logout", h.Logout)
}

func setSessionCookie(w http.ResponseWriter, name, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
