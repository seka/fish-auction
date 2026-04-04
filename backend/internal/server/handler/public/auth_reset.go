package public

import (
	"encoding/json"
	"errors"
	"net/http"

	domainErrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
)

// AuthResetHandler handles HTTP requests related to password resets.
type AuthResetHandler struct {
	reg registry.UseCase
}

// NewAuthResetHandler creates a new AuthResetHandler instance.
func NewAuthResetHandler(r registry.UseCase) *AuthResetHandler {
	return &AuthResetHandler{reg: r}
}

// RequestReset handles the password reset request.
func (h *AuthResetHandler) RequestReset(w http.ResponseWriter, r *http.Request) {
	var req request.ResetPassword
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	uc := h.reg.NewRequestPasswordResetUseCase()
	// Process reset request
	if err := uc.Execute(r.Context(), req.Email); err != nil {
		var notFoundErr *domainErrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			// Security: Don't reveal if user exists.
			// Return 200 OK even if not found.
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response.Message{Message: "If the email exists, a reset link has been sent."})
			return
		}
		// System errors (DB, Email, etc.) are 500
		util.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response.Message{Message: "If the email exists, a reset link has been sent."})
}

// VerifyToken provides VerifyToken related functionality.
func (h *AuthResetHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	var req request.ResetPasswordVerify
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	uc := h.reg.NewVerifyResetTokenUseCase()
	if err := uc.Execute(r.Context(), req.Token); err != nil {
		var unauthErr *domainErrors.UnauthorizedError
		if errors.As(err, &unauthErr) {
			util.WriteError(w, http.StatusUnauthorized, unauthErr.Message)
			return
		}
		util.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response.Message{Message: "Token is valid"})
}

// ConfirmReset provides ConfirmReset related functionality.
func (h *AuthResetHandler) ConfirmReset(w http.ResponseWriter, r *http.Request) {
	var req request.ResetPasswordConfirm
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	uc := h.reg.NewResetPasswordUseCase()
	if err := uc.Execute(r.Context(), req.Token, req.NewPassword); err != nil {
		var notFoundErr *domainErrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			util.WriteError(w, http.StatusBadRequest, "Invalid or expired token")
			return
		}

		var unauthErr *domainErrors.UnauthorizedError
		if errors.As(err, &unauthErr) {
			util.WriteError(w, http.StatusBadRequest, unauthErr.Message)
			return
		}

		var valErr *domainErrors.ValidationError
		if errors.As(err, &valErr) {
			util.WriteError(w, http.StatusBadRequest, valErr.Message)
			return
		}

		util.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response.Message{Message: "Password updated successfully"})
}

// RegisterRoutes registers the password reset handler routes to the given mux.
func (h *AuthResetHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/password-reset/request", h.RequestReset)
	mux.HandleFunc("POST /auth/password-reset/verify", h.VerifyToken)
	mux.HandleFunc("POST /auth/password-reset/confirm", h.ConfirmReset)
}
