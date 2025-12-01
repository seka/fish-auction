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

type VenueHandler struct {
	createUseCase venue.CreateVenueUseCase
	listUseCase   venue.ListVenuesUseCase
	getUseCase    venue.GetVenueUseCase
	updateUseCase venue.UpdateVenueUseCase
	deleteUseCase venue.DeleteVenueUseCase
}

func NewVenueHandler(r registry.UseCase) *VenueHandler {
	return &VenueHandler{
		createUseCase: r.NewCreateVenueUseCase(),
		listUseCase:   r.NewListVenuesUseCase(),
		getUseCase:    r.NewGetVenueUseCase(),
		updateUseCase: r.NewUpdateVenueUseCase(),
		deleteUseCase: r.NewDeleteVenueUseCase(),
	}
}

func (h *VenueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateVenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	venue := &model.Venue{
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}

	created, err := h.createUseCase.Execute(r.Context(), venue)
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
	json.NewEncoder(w).Encode(resp)
}

func (h *VenueHandler) List(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(resp)
}

func (h *VenueHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/venues/"):]
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
	json.NewEncoder(w).Encode(resp)
}

func (h *VenueHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/venues/"):]
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

	venue := &model.Venue{
		ID:          id,
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}

	if err := h.updateUseCase.Execute(r.Context(), venue); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *VenueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/venues/"):]
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

func (h *VenueHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/venues", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			h.List(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/venues/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Get(w, r)
		case http.MethodPut:
			h.Update(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
