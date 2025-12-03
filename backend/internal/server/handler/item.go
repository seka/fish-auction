package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
)

type ItemHandler struct {
	createUseCase item.CreateItemUseCase
	listUseCase   item.ListItemsUseCase
}

func NewItemHandler(r registry.UseCase) *ItemHandler {
	return &ItemHandler{
		createUseCase: r.NewCreateItemUseCase(),
		listUseCase:   r.NewListItemsUseCase(),
	}
}

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

	resp := dto.ItemResponse{
		ID:                created.ID,
		AuctionID:         created.AuctionID,
		FishermanID:       created.FishermanID,
		FishType:          created.FishType,
		Quantity:          created.Quantity,
		Unit:              created.Unit,
		Status:            created.Status.String(),
		HighestBid:        created.HighestBid,
		HighestBidderID:   created.HighestBidderID,
		HighestBidderName: created.HighestBidderName,
		CreatedAt:         created.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	items, err := h.listUseCase.Execute(r.Context(), status)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		resp[i] = dto.ItemResponse{
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
			CreatedAt:         item.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ItemHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
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
