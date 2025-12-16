package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"
)

type AdminHandler struct {
	updatePasswordUseCase admin.UpdatePasswordUseCase
}

func NewAdminHandler(r registry.UseCase) *AdminHandler {
	return &AdminHandler{
		updatePasswordUseCase: r.NewAdminUpdatePasswordUseCase(),
	}
}

func (h *AdminHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	cookie, err := r.Cookie("admin_session")
	if err != nil || cookie.Value != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get admin ID from cookie
	idCookie, err := r.Cookie("admin_id")
	if err != nil {
		http.Error(w, "Unauthorized: Admin ID missing", http.StatusUnauthorized)
		return
	}

	var adminID int
	if _, err := fmt.Sscanf(idCookie.Value, "%d", &adminID); err != nil {
		http.Error(w, "Invalid admin ID", http.StatusBadRequest)
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
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/admin/password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			h.UpdatePassword(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
