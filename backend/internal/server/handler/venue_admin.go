package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

// AdminVenueHandler handles admin HTTP requests related to venues.
type AdminVenueHandler struct {
	createUseCase venue.CreateVenueUseCase
	updateUseCase venue.UpdateVenueUseCase
	deleteUseCase venue.DeleteVenueUseCase
}

// NewAdminVenueHandler creates a new AdminVenueHandler instance.
func NewAdminVenueHandler(r registry.UseCase) *AdminVenueHandler {
	return &AdminVenueHandler{
		createUseCase: r.NewCreateVenueUseCase(),
		updateUseCase: r.NewUpdateVenueUseCase(),
		deleteUseCase: r.NewDeleteVenueUseCase(),
	}
}

// Create handles the venue creation request.
func (h *AdminVenueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	vn := &model.Venue{
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}

	created, err := h.createUseCase.Execute(r.Context(), vn)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := dto.VenueResponse{
		ID:          created.ID,
		Name:        created.Name,
		Location:    created.Location,
		Description: created.Description,
		CreatedAt:   created.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Update handles the request to update a specific venue.
func (h *AdminVenueHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	vn := &model.Venue{
		ID:          id,
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}

	if err := h.updateUseCase.Execute(r.Context(), vn); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete handles the venue deletion request.
func (h *AdminVenueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes registers the admin venue handler routes to the given mux.
func (h *AdminVenueHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /venues", h.Create)
	mux.HandleFunc("PUT /venues/{id}", h.Update)
	mux.HandleFunc("DELETE /venues/{id}", h.Delete)
}
