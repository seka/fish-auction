package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"strconv"

	"github.com/gorilla/mux"
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

	item := &model.AuctionItem{
		AuctionID:   req.AuctionID,
		FishermanID: req.FishermanID,
		FishType:    req.FishType,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
	}

	created, err := h.createUseCase.Execute(r.Context(), item)
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
	for i, item := range items {
		resp[i] = h.toResponse(&item)
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

// Update handles the request to update a specific item.
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path: /items/{id} or /api/admin/items/{id}
	path := r.URL.Path
	segments := strings.Split(strings.Trim(path, "/"), "/")
	var idStr string
	for i, s := range segments {
		if s == "items" && i+1 < len(segments) {
			idStr = segments[i+1]
			break
		}
	}

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
	// Extract ID from path: /items/{id}
	path := r.URL.Path
	segments := strings.Split(strings.Trim(path, "/"), "/")
	var idStr string
	for i, s := range segments {
		if s == "items" && i+1 < len(segments) {
			idStr = segments[i+1]
			break
		}
	}

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
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
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
	vars := mux.Vars(r)
	auctionID, err := strconv.Atoi(vars["id"])
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

func (h *ItemHandler) toResponse(item *model.AuctionItem) dto.ItemResponse {
	return dto.ItemResponse{
		ID:                item.ID,
		AuctionID:         item.AuctionID,
		FishermanID:       item.FishermanID,
		FishType:          item.FishType,
		Quantity:          item.Quantity,
		Unit:              item.Unit,
		Status:            item.Status.String(),
		HighestBid:        item.HighestBid,
		HighestBidderID:   item.HighestBidderID,
		HighestBidderName: item.HighestBidderName,
		SortOrder:         item.SortOrder,
		CreatedAt:         item.CreatedAt,
	}
}

// RegisterRoutes registers the item handler routes to the given router.
func (h *ItemHandler) RegisterRoutes(r *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	r.HandleFunc("/api/items", h.List).Methods(http.MethodGet)
	r.Handle("/api/items", authMiddleware(http.HandlerFunc(h.Create))).Methods(http.MethodPost)
	r.Handle("/api/items/{id:[0-9]+}", authMiddleware(http.HandlerFunc(h.Update))).Methods(http.MethodPut)
	r.Handle("/api/items/{id:[0-9]+}", authMiddleware(http.HandlerFunc(h.Delete))).Methods(http.MethodDelete)
	r.Handle("/api/items/{id:[0-9]+}/sort-order", authMiddleware(http.HandlerFunc(h.UpdateSortOrder))).Methods(http.MethodPut)
}
