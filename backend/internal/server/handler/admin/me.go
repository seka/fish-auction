package admin

import (
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
	"github.com/seka/fish-auction/backend/internal/server/util"
)

// MeHandler returns the currently authenticated admin's profile.
type MeHandler struct {
	adminRepo repository.AdminRepository
}

// NewMeHandler creates a new MeHandler instance.
func NewMeHandler(adminRepo repository.AdminRepository) *MeHandler {
	return &MeHandler{adminRepo: adminRepo}
}

// GetMe returns the admin profile corresponding to the session.
func (h *MeHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.AdminIDFromContext(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	admin, err := h.adminRepo.FindByID(r.Context(), adminID)
	if err != nil || admin == nil {
		util.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	util.WriteJSON(w, http.StatusOK, struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{ID: admin.ID, Email: admin.Email})
}

// RegisterRoutes registers the me handler routes to the given mux.
func (h *MeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /me", h.GetMe)
}
