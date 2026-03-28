package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
)

// AdminAuthResetHandler handles HTTP requests related to admin password resets.
type AdminAuthResetHandler struct {
	reg registry.UseCase
}

// NewAdminAuthResetHandler creates a new AdminAuthResetHandler instance.
func NewAdminAuthResetHandler(r registry.UseCase) *AdminAuthResetHandler {
	return &AdminAuthResetHandler{
		reg: r,
	}
}

// RequestReset handles the admin password reset request.
func (h *AdminAuthResetHandler) RequestReset(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	uc := h.reg.NewRequestAdminPasswordResetUseCase()
	if err := uc.Execute(r.Context(), req.Email); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Password reset email sent if account exists"})
}

// VerifyToken provides VerifyToken related functionality.
func (h *AdminAuthResetHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	uc := h.reg.NewVerifyAdminResetTokenUseCase()
	if err := uc.Execute(r.Context(), req.Token); err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Token is valid"})
}

// ConfirmReset provides ConfirmReset related functionality.
func (h *AdminAuthResetHandler) ConfirmReset(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	uc := h.reg.NewResetAdminPasswordUseCase()
	if err := uc.Execute(r.Context(), req.Token, req.NewPassword); err != nil {
		if err.Error() == "Invalid or expired token" || err.Error() == "token expired" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

// RegisterRoutes registers the admin password reset handler routes to the given mux.
func (h *AdminAuthResetHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/admin/password-reset/request", h.RequestReset)
	mux.HandleFunc("POST /api/admin/password-reset/verify", h.VerifyToken)
	mux.HandleFunc("POST /api/admin/password-reset/confirm", h.ConfirmReset)
}
