package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

// AdminBuyerHandler handles admin HTTP requests related to buyers.
type AdminBuyerHandler struct {
	createUseCase buyer.CreateBuyerUseCase
	listUseCase   buyer.ListBuyersUseCase
	deleteUseCase buyer.DeleteBuyerUseCase
}

// NewAdminBuyerHandler creates a new AdminBuyerHandler instance.
func NewAdminBuyerHandler(r registry.UseCase) *AdminBuyerHandler {
	return &AdminBuyerHandler{
		createUseCase: r.NewCreateBuyerUseCase(),
		listUseCase:   r.NewListBuyersUseCase(),
		deleteUseCase: r.NewDeleteBuyerUseCase(),
	}
}

// Create handles the buyer creation request.
func (h *AdminBuyerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBuyerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}
	buy, err := h.createUseCase.Execute(r.Context(), req.Name, req.Email, req.Password, req.Organization, req.ContactInfo)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	resp := dto.BuyerResponse{ID: buy.ID, Name: buy.Name}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// List handles the request to list buyers.
func (h *AdminBuyerHandler) List(w http.ResponseWriter, r *http.Request) {
	buyers, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}
	resp := make([]dto.BuyerResponse, len(buyers))
	for i, b := range buyers {
		resp[i] = dto.BuyerResponse{ID: b.ID, Name: b.Name}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Delete handles the buyer deletion request.
func (h *AdminBuyerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid buyer id")
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes registers the admin buyer handler routes to the given mux.
func (h *AdminBuyerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /buyers", h.List)
	mux.HandleFunc("POST /buyers", h.Create)
	mux.HandleFunc("DELETE /buyers/{id}", h.Delete)
}
