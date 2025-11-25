package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
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
		util.HandleError(w, err)
		return
	}

	bid := &model.Bid{
		ItemID:  req.ItemID,
		BuyerID: req.BuyerID,
		Price:   req.Price,
	}

	if _, err := h.useCase.Bid(r.Context(), bid); err != nil {
		util.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bid placed successfully"})
}

func (h *BidHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/bids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
