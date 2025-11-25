package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/usecase"
)

type FishermanHandler struct {
	useCase usecase.FishermanUseCase
}

func NewFishermanHandler(uc usecase.FishermanUseCase) *FishermanHandler {
	return &FishermanHandler{useCase: uc}
}

func (h *FishermanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fisherman, err := h.useCase.Create(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fisherman)
}

func (h *FishermanHandler) List(w http.ResponseWriter, r *http.Request) {
	fishermen, err := h.useCase.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fishermen)
}

func (h *FishermanHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/fishermen", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			h.List(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
