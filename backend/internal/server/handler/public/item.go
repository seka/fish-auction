package public

import (
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
)

// ItemHandler handles public HTTP requests related to items.
type ItemHandler struct {
	listUseCase item.ListItemsUseCase
}

// NewItemHandler creates a new ItemHandler instance.
func NewItemHandler(r registry.UseCase) *ItemHandler {
	return &ItemHandler{
		listUseCase: r.NewListItemsUseCase(),
	}
}

// List handles the request to list items.
func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]response.Item, len(items))
	for i := range items {
		it := items[i]
		resp[i] = h.toResponse(&it)
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

func (h *ItemHandler) toResponse(it *model.AuctionItem) response.Item {
	var highestBid *int
	if it.HighestBid != nil {
		amt := it.HighestBid.Amount()
		highestBid = &amt
	}
	return response.Item{
		ID:          it.ID,
		AuctionID:   it.AuctionID,
		FishermanID: it.FishermanID,
		FishType:    it.FishType,
		Quantity:    it.Quantity,
		Unit:        it.Unit,
		HighestBid:  highestBid,
		SortOrder:   it.SortOrder,
		CreatedAt:   it.CreatedAt,
	}
}

// RegisterRoutes registers the public item handler routes to the given mux.
func (h *ItemHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/items", h.List)
}
