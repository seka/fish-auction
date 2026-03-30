package public

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/venue"
)

// VenueHandler handles public HTTP requests related to venues.
type VenueHandler struct {
	listUseCase venue.ListVenuesUseCase
	getUseCase  venue.GetVenueUseCase
}

// NewVenueHandler creates a new VenueHandler instance.
func NewVenueHandler(r registry.UseCase) *VenueHandler {
	return &VenueHandler{
		listUseCase: r.NewListVenuesUseCase(),
		getUseCase:  r.NewGetVenueUseCase(),
	}
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
		resp[i] = response.Venue{
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

// Get handles the request to get a single venue.
func (h *VenueHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid venue ID", http.StatusBadRequest)
		return
	}

	v, err := h.getUseCase.Execute(r.Context(), id)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := response.Venue{
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
func (h *VenueHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/venues", h.List)
	mux.HandleFunc("GET /api/venues/{id}", h.Get)
}
