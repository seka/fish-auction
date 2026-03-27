package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

// PublicVenueHandler handles public HTTP requests related to venues.
type PublicVenueHandler struct {
	listUseCase venue.ListVenuesUseCase
	getUseCase  venue.GetVenueUseCase
}

// NewPublicVenueHandler creates a new PublicVenueHandler instance.
func NewPublicVenueHandler(r registry.UseCase) *PublicVenueHandler {
	return &PublicVenueHandler{
		listUseCase: r.NewListVenuesUseCase(),
		getUseCase:  r.NewGetVenueUseCase(),
	}
}

// List handles the request to list venues.
func (h *PublicVenueHandler) List(w http.ResponseWriter, r *http.Request) {
	venues, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]dto.VenueResponse, len(venues))
	for i, v := range venues {
		resp[i] = dto.VenueResponse{
			ID:          v.ID,
			Name:        v.Name,
			Location:    v.Location,
			Description: v.Description,
			CreatedAt:   v.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Get handles the request to get a specific venue.
func (h *PublicVenueHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	v, err := h.getUseCase.Execute(r.Context(), id)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := dto.VenueResponse{
		ID:          v.ID,
		Name:        v.Name,
		Location:    v.Location,
		Description: v.Description,
		CreatedAt:   v.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// RegisterRoutes registers the public venue handler routes to the given mux.
func (h *PublicVenueHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/venues", h.List)
	mux.HandleFunc("GET /api/venues/{id}", h.Get)
}
