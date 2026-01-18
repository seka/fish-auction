package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/fisherman"
)

type FishermanHandler struct {
	createUseCase fisherman.CreateFishermanUseCase
	listUseCase   fisherman.ListFishermenUseCase
	deleteUseCase fisherman.DeleteFishermanUseCase
}

func NewFishermanHandler(r registry.UseCase) *FishermanHandler {
	return &FishermanHandler{
		createUseCase: r.NewCreateFishermanUseCase(),
		listUseCase:   r.NewListFishermenUseCase(),
		deleteUseCase: r.NewDeleteFishermanUseCase(),
	}
}

func (h *FishermanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFishermanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	fisherman, err := h.createUseCase.Execute(r.Context(), req.Name)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := dto.FishermanResponse{ID: fisherman.ID, Name: fisherman.Name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

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
	json.NewEncoder(w).Encode(resp)
}

func (h *FishermanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := util.ParseID(r)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FishermanHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/fishermen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.List(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
