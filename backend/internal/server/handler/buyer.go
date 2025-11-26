package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

type BuyerHandler struct {
	createUseCase buyer.CreateBuyerUseCase
	listUseCase   buyer.ListBuyersUseCase
}

func NewBuyerHandler(r registry.UseCase) *BuyerHandler {
	return &BuyerHandler{
		createUseCase: r.NewCreateBuyerUseCase(),
		listUseCase:   r.NewListBuyersUseCase(),
	}
}

func (h *BuyerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBuyerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}
	buyer, err := h.createUseCase.Execute(r.Context(), req.Name)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	resp := dto.BuyerResponse{ID: buyer.ID, Name: buyer.Name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) List(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(resp)
}

func (h *BuyerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/buyers", func(w http.ResponseWriter, r *http.Request) {
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
