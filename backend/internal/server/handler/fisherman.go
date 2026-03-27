package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/fisherman"
)

// FishermanHandler handles HTTP requests related to fishermen.
type FishermanHandler struct {
	createUseCase fisherman.CreateFishermanUseCase
	listUseCase   fisherman.ListFishermenUseCase
	deleteUseCase fisherman.DeleteFishermanUseCase
}

// NewFishermanHandler creates a new FishermanHandler instance.
func NewFishermanHandler(r registry.UseCase) *FishermanHandler {
	return &FishermanHandler{
		createUseCase: r.NewCreateFishermanUseCase(),
		listUseCase:   r.NewListFishermenUseCase(),
		deleteUseCase: r.NewDeleteFishermanUseCase(),
	}
}

// Create handles the fisherman creation request.
func (h *FishermanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFishermanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	fm, err := h.createUseCase.Execute(r.Context(), req.Name)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := dto.FishermanResponse{ID: fm.ID, Name: fm.Name}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// List handles the request to list fishermen.
func (h *FishermanHandler) List(w http.ResponseWriter, r *http.Request) {
	fishermen, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]dto.FishermanResponse, len(fishermen))
	for i, f := range fishermen {
		resp[i] = dto.FishermanResponse{ID: f.ID, Name: f.Name}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Delete handles the fisherman deletion request.
func (h *FishermanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid fisherman id")
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RegisterRoutes registers the fisherman handler routes to the given mux.
func (h *FishermanHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/fishermen", h.List)
	mux.HandleFunc("POST /api/fishermen", h.Create)
	mux.HandleFunc("DELETE /api/fishermen/{id}", h.Delete)
}
