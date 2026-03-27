package handler

import (
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
)

// PublicItemHandler handles public HTTP requests related to items.
type PublicItemHandler struct {
	listUseCase item.ListItemsUseCase
}

// NewPublicItemHandler creates a new PublicItemHandler instance.
func NewPublicItemHandler(r registry.UseCase) *PublicItemHandler {
	return &PublicItemHandler{
		listUseCase: r.NewListItemsUseCase(),
	}
}

// List handles the request to list items.
func (h *PublicItemHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	items, err := h.listUseCase.Execute(r.Context(), status)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]dto.ItemResponse, len(items))
	for i := range items {
		it := items[i]
		resp[i] = h.toResponse(&it)
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

func (h *PublicItemHandler) toResponse(it *model.AuctionItem) dto.ItemResponse {
	var highestBid *int
	if it.HighestBid != nil {
		amt := it.HighestBid.Amount()
		highestBid = &amt
	}
	return dto.ItemResponse{
		ID:                it.ID,
		AuctionID:         it.AuctionID,
		FishermanID:       it.FishermanID,
		FishType:          it.FishType,
		Quantity:          it.Quantity,
		Unit:              it.Unit,
		Status:            it.Status.String(),
		HighestBid:        highestBid,
		HighestBidderID:   it.HighestBidderID,
		HighestBidderName: it.HighestBidderName,
		SortOrder:         it.SortOrder,
		CreatedAt:         it.CreatedAt,
	}
}

// RegisterRoutes registers the public item handler routes to the given mux.
func (h *PublicItemHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/items", h.List)
}
