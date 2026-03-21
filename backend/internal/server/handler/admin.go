package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
)

// AdminHandler handles HTTP requests related to administration.
type AdminHandler struct {
	updatePasswordUseCase admin.UpdatePasswordUseCase
}

// NewAdminHandler creates a new AdminHandler instance.
func NewAdminHandler(r registry.UseCase) *AdminHandler {
	return &AdminHandler{
		updatePasswordUseCase: r.NewAdminUpdatePasswordUseCase(),
	}
}

// UpdatePassword handles the password update request for an admin.
func (h *AdminHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.AdminIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	if err := h.updatePasswordUseCase.Execute(r.Context(), adminID, req.CurrentPassword, req.NewPassword); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dto.MessageResponse{Message: "Password updated successfully"}); err != nil {
		util.HandleError(w, err)
		return
	}
}

// RegisterRoutes registers the admin handler routes to the given mux.
func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/admin/password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			h.UpdatePassword(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
