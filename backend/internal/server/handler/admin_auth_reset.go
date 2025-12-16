package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
)

type AdminAuthResetHandler struct {
	requestUseCase admin.RequestPasswordResetUseCase
	resetUseCase   admin.ResetPasswordUseCase
}

func NewAdminAuthResetHandler(uc registry.UseCase) *AdminAuthResetHandler {
	return &AdminAuthResetHandler{
		requestUseCase: uc.NewRequestAdminPasswordResetUseCase(),
		resetUseCase:   uc.NewResetAdminPasswordUseCase(),
	}
}

func (h *AdminAuthResetHandler) RequestReset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.requestUseCase.Execute(r.Context(), req.Email); err != nil {
		// Log error but generally return success to prevent email enumeration
		// Or return 500 if it's a system error?
		// For now, simple error response if something major fails, or success.
		// UseCase returns nil if user not found, so this error would be infrastructure.
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset email sent if account exists"})
}

func (h *AdminAuthResetHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	// Optional: verify token validity without resetting password yet
	// Not strictly required if frontend just prompts for new password.
	// But nice for UX (show "invalid token" immediately).
	// Current UseCase doesn't support separate verify, so we skip or enhance UseCase.
	// For now, skipped to keep simple. Frontend can just let user submit password.
	w.WriteHeader(http.StatusOK)
}

func (h *AdminAuthResetHandler) ConfirmReset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.resetUseCase.Execute(r.Context(), req.Token, req.NewPassword); err != nil {
		if err.Error() == "invalid or expired token" || err.Error() == "token expired" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

func (h *AdminAuthResetHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/admin/password-reset/request", h.RequestReset)
	mux.HandleFunc("POST /api/admin/password-reset/verify", h.VerifyToken)
	mux.HandleFunc("POST /api/admin/password-reset/confirm", h.ConfirmReset)
}
