package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/bid"
	"github.com/seka/fish-auction/backend/internal/usecase/invoice"
)

type BidHandler struct {
	createUseCase       bid.CreateBidUseCase
	listInvoicesUseCase invoice.ListInvoicesUseCase
}

func NewBidHandler(r registry.UseCase) *BidHandler {
	return &BidHandler{
		createUseCase:       r.NewCreateBidUseCase(),
		listInvoicesUseCase: r.NewListInvoicesUseCase(),
	}
}

func (h *BidHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	buyerID, ok := r.Context().Value("buyer_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get item to find auction_id
	// Since we don't have GetItemUseCase, we'll need to get auction_id from the item
	// For now, we'll skip item lookup and assume the client provides correct auction context
	// In production, you should add GetItemUseCase or use ItemRepository directly

	// For MVP, we'll validate auction period based on item's auction_id
	// This requires adding auction_id to the bid request or looking it up
	// For now, let's comment out the period check and implement it properly later

	bid := &model.Bid{
		ItemID:  req.ItemID,
		BuyerID: buyerID,
		Price:   req.Price,
	}

	if _, err := h.createUseCase.Execute(r.Context(), bid); err != nil {
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
