package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
)

// ItemHandler handles HTTP requests related to items.
type ItemHandler struct {
	createUseCase          item.CreateItemUseCase
	listUseCase            item.ListItemsUseCase
	updateUseCase          item.UpdateItemUseCase
	deleteUseCase          item.DeleteItemUseCase
	updateSortOrderUseCase item.UpdateItemSortOrderUseCase
	reorderItemsUseCase    item.ReorderItemsUseCase
}

// NewItemHandler creates a new ItemHandler instance.
func NewItemHandler(r registry.UseCase) *ItemHandler {
	return &ItemHandler{
		createUseCase:          r.NewCreateItemUseCase(),
		listUseCase:            r.NewListItemsUseCase(),
		updateUseCase:          r.NewUpdateItemUseCase(),
		deleteUseCase:          r.NewDeleteItemUseCase(),
		updateSortOrderUseCase: r.NewUpdateItemSortOrderUseCase(),
		reorderItemsUseCase:    r.NewReorderItemsUseCase(),
	}
}

// Create handles the item creation request.
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateItemRequest
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

// List handles the request to list items.
func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
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

// Update handles the request to update a specific item.
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req dto.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	itemModel := &model.AuctionItem{
		ID:          id,
		AuctionID:   req.AuctionID,
		FishermanID: req.FishermanID,
		FishType:    req.FishType,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		Status:      model.ItemStatus(req.Status),
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
		util.WriteError(w, http.StatusBadRequest, "invalid item id")
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
		util.WriteError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req dto.UpdateItemSortOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.updateSortOrderUseCase.Execute(r.Context(), id, req.SortOrder); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Reorder handles the request to reorder items within an auction.
func (h *ItemHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	auctionID, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid auction id")
		return
	}

	var req dto.ReorderItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.reorderItemsUseCase.Execute(r.Context(), auctionID, req.IDs); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ItemHandler) toResponse(it *model.AuctionItem) dto.ItemResponse {
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

// RegisterRoutes registers the item handler routes to the given mux.
func (h *ItemHandler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /api/items", h.List)
	mux.Handle("POST /api/items", authMiddleware(http.HandlerFunc(h.Create)))
	mux.Handle("PUT /api/items/{id}", authMiddleware(http.HandlerFunc(h.Update)))
	mux.Handle("DELETE /api/items/{id}", authMiddleware(http.HandlerFunc(h.Delete)))
	mux.Handle("PUT /api/items/{id}/sort-order", authMiddleware(http.HandlerFunc(h.UpdateSortOrder)))
}
