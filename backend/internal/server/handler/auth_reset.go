package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
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
	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	uc := h.reg.NewRequestPasswordResetUseCase()
	// Process reset request
	_ = uc.Execute(r.Context(), req.Email)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "If the email exists, a reset link has been sent."})
}

// VerifyToken provides VerifyToken related functionality.
func (h *AuthResetHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	uc := h.reg.NewVerifyResetTokenUseCase()
	if err := uc.Execute(r.Context(), req.Token); err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Token is valid"})
}

// ConfirmReset provides ConfirmReset related functionality.
func (h *AuthResetHandler) ConfirmReset(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	uc := h.reg.NewResetPasswordUseCase()
	if err := uc.Execute(r.Context(), req.Token, req.NewPassword); err != nil {
		if err.Error() == "Invalid or expired token" || err.Error() == "token expired" {
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

// RegisterRoutes registers the password reset handler routes to the given mux.
func (h *AuthResetHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/password-reset/request", h.RequestReset)
	mux.HandleFunc("POST /api/auth/password-reset/verify", h.VerifyToken)
	mux.HandleFunc("POST /api/auth/password-reset/confirm", h.ConfirmReset)
}
