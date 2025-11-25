package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/usecase"
)

type BidHandler struct {
	useCase usecase.BidUseCase
}

func NewBidHandler(uc usecase.BidUseCase) *BidHandler {
	return &BidHandler{useCase: uc}
}

func (h *BidHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction := &model.Transaction{
		ItemID:  req.ItemID,
		BuyerID: req.BuyerID,
		Price:   req.Price,
	}

	result, err := h.useCase.Bid(r.Context(), transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.BidResponse{
		ID:        result.ID,
		ItemID:    result.ItemID,
		BuyerID:   result.BuyerID,
		Price:     result.Price,
		CreatedAt: result.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *BidHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/bid", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
