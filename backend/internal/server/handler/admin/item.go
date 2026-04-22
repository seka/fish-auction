package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
)

// ItemHandler handles admin HTTP requests related to items.
type ItemHandler struct {
	createUseCase          item.CreateItemUseCase
	updateUseCase          item.UpdateItemUseCase
	deleteUseCase          item.DeleteItemUseCase
	updateSortOrderUseCase item.UpdateItemSortOrderUseCase
}

// NewItemHandler creates a new ItemHandler instance.
func NewItemHandler(r registry.UseCase) *ItemHandler {
	return &ItemHandler{
		createUseCase:          r.NewCreateItemUseCase(),
		updateUseCase:          r.NewUpdateItemUseCase(),
		deleteUseCase:          r.NewDeleteItemUseCase(),
		updateSortOrderUseCase: r.NewUpdateItemSortOrderUseCase(),
	}
}

// Create handles the item creation request.
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateItem
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	it := &model.AuctionItem{
		AuctionID:   req.AuctionID,
		FishermanID: req.FishermanID,
		FishType:    req.FishType,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
	}

	created, err := h.createUseCase.Execute(r.Context(), it)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, h.toResponse(created))
}

// Update handles the request to update a specific item.
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req request.UpdateItem
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	itemModel := &model.AuctionItem{
		ID:          id,
		AuctionID:   req.AuctionID,
		FishermanID: req.FishermanID,
		FishType:    req.FishType,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
	}

	updated, err := h.updateUseCase.Execute(r.Context(), itemModel)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, h.toResponse(updated))
}

// Delete handles the item deletion request.
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateSortOrder handles the request to update an item's sort order.
func (h *ItemHandler) UpdateSortOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req request.UpdateItemSortOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.updateSortOrderUseCase.Execute(r.Context(), id, req.SortOrder); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ItemHandler) toResponse(it *model.AuctionItem) response.Item {
	var highestBid *int
	if it.HighestBid != nil {
		amt := it.HighestBid.Amount()
		highestBid = &amt
	}
	return response.Item{
		ID:                it.ID,
		AuctionID:         it.AuctionID,
		FishermanID:       it.FishermanID,
		FishType:          it.FishType,
		Quantity:          it.Quantity,
		Unit:              it.Unit,
		HighestBid:        highestBid,
		HighestBidderID:   it.HighestBidderID,
		HighestBidderName: it.HighestBidderName,
		SortOrder:         it.SortOrder,
		CreatedAt:         it.CreatedAt,
	}
}

// RegisterRoutes registers the admin item handler routes to the given mux.
func (h *ItemHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /items", h.Create)
	mux.HandleFunc("PUT /items/{id}", h.Update)
	mux.HandleFunc("DELETE /items/{id}", h.Delete)
	mux.HandleFunc("PUT /items/{id}/sort-order", h.UpdateSortOrder)
}
