package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

// VenueHandler handles admin HTTP requests related to venues.
type VenueHandler struct {
	createUseCase venue.CreateVenueUseCase
	listUseCase   venue.ListVenuesUseCase
	updateUseCase venue.UpdateVenueUseCase
	deleteUseCase venue.DeleteVenueUseCase
}

// NewVenueHandler creates a new VenueHandler instance.
func NewVenueHandler(r registry.UseCase) *VenueHandler {
	return &VenueHandler{
		createUseCase: r.NewCreateVenueUseCase(),
		listUseCase:   r.NewListVenuesUseCase(),
		updateUseCase: r.NewUpdateVenueUseCase(),
		deleteUseCase: r.NewDeleteVenueUseCase(),
	}
}

// Create handles the venue creation request.
func (h *VenueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateVenue
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

	util.WriteJSON(w, http.StatusCreated, h.toResponse(*created))
}

// List handles the request to list venues.
func (h *VenueHandler) List(w http.ResponseWriter, r *http.Request) {
	venues, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]response.Venue, len(venues))
	for i, v := range venues {
		resp[i] = h.toResponse(v)
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

// Update handles the request to update a specific venue.
func (h *VenueHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid venue id")
		return
	}

	var req request.UpdateVenue
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
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
func (h *VenueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid venue id")
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *VenueHandler) toResponse(v model.Venue) response.Venue {
	return response.Venue{
		ID:          v.ID,
		Name:        v.Name,
		Location:    v.Location,
		Description: v.Description,
		CreatedAt:   v.CreatedAt,
	}
}

// RegisterRoutes registers the admin venue handler routes to the given mux.
func (h *VenueHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /venues", h.Create)
	mux.HandleFunc("GET /venues", h.List)
	mux.HandleFunc("PUT /venues/{id}", h.Update)
	mux.HandleFunc("DELETE /venues/{id}", h.Delete)
}
