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

// AdminItemHandler handles admin HTTP requests related to items.
type AdminItemHandler struct {
	createUseCase          item.CreateItemUseCase
	updateUseCase          item.UpdateItemUseCase
	deleteUseCase          item.DeleteItemUseCase
	updateSortOrderUseCase item.UpdateItemSortOrderUseCase
}

// NewAdminItemHandler creates a new AdminItemHandler instance.
func NewAdminItemHandler(r registry.UseCase) *AdminItemHandler {
	return &AdminItemHandler{
		createUseCase:          r.NewCreateItemUseCase(),
		updateUseCase:          r.NewUpdateItemUseCase(),
		deleteUseCase:          r.NewDeleteItemUseCase(),
		updateSortOrderUseCase: r.NewUpdateItemSortOrderUseCase(),
	}
}

// Create handles the item creation request.
func (h *AdminItemHandler) Create(w http.ResponseWriter, r *http.Request) {
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

// Update handles the request to update a specific item.
func (h *AdminItemHandler) Update(w http.ResponseWriter, r *http.Request) {
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
func (h *AdminItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
func (h *AdminItemHandler) UpdateSortOrder(w http.ResponseWriter, r *http.Request) {
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


func (h *AdminItemHandler) toResponse(it *model.AuctionItem) dto.ItemResponse {
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

// RegisterRoutes registers the admin item handler routes to the given mux.
func (h *AdminItemHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /items", h.Create)
	mux.HandleFunc("PUT /items/{id}", h.Update)
	mux.HandleFunc("DELETE /items/{id}", h.Delete)
	mux.HandleFunc("PUT /items/{id}/sort-order", h.UpdateSortOrder)
}
