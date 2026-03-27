package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

// PublicAuctionHandler handles public HTTP requests related to auctions.
type PublicAuctionHandler struct {
	listUseCase     auction.ListAuctionsUseCase
	getUseCase      auction.GetAuctionUseCase
	getItemsUseCase auction.GetAuctionItemsUseCase
}

// NewPublicAuctionHandler creates a new PublicAuctionHandler instance.
func NewPublicAuctionHandler(r registry.UseCase) *PublicAuctionHandler {
	return &PublicAuctionHandler{
		listUseCase:     r.NewListAuctionsUseCase(),
		getUseCase:      r.NewGetAuctionUseCase(),
		getItemsUseCase: r.NewGetAuctionItemsUseCase(),
	}
}

// List handles the request to list auctions.
func (h *PublicAuctionHandler) List(w http.ResponseWriter, r *http.Request) {
	filters := &repository.AuctionFilters{}

	if venueIDStr := r.URL.Query().Get("venue_id"); venueIDStr != "" {
		if id, err := strconv.Atoi(venueIDStr); err == nil {
			filters.VenueID = &id
		}
	}

	if dateStr := r.URL.Query().Get("date"); dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			filters.AuctionDate = &date
		}
	}

	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		status := model.AuctionStatus(statusStr)
		filters.Status = &status
	}

	auctions, err := h.listUseCase.Execute(r.Context(), filters)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]dto.AuctionResponse, len(auctions))
	for i, a := range auctions {
		resp[i] = dto.AuctionResponse{
			ID:          a.ID,
			VenueID:     a.VenueID,
			AuctionDate: a.Period.AuctionDate.Format("2006-01-02"),
			StartTime:   formatTime(a.Period.StartAt),
			EndTime:     formatTime(a.Period.EndAt),
			Status:      string(a.Status),
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Get handles the request to get a specific auction.
func (h *PublicAuctionHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	a, err := h.getUseCase.Execute(r.Context(), id)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := dto.AuctionResponse{
		ID:          a.ID,
		VenueID:     a.VenueID,
		AuctionDate: a.Period.AuctionDate.Format("2006-01-02"),
		StartTime:   formatTime(a.Period.StartAt),
		EndTime:     formatTime(a.Period.EndAt),
		Status:      string(a.Status),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// GetItems handles the request to get items for a specific auction.
func (h *PublicAuctionHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	items, err := h.getItemsUseCase.Execute(r.Context(), id)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]dto.ItemResponse, len(items))
	for i := range items {
		item := items[i]
		var highestBid *int
		if item.HighestBid != nil {
			amt := item.HighestBid.Amount()
			highestBid = &amt
		}
		resp[i] = dto.ItemResponse{
			ID:                item.ID,
			AuctionID:         item.AuctionID,
			FishermanID:       item.FishermanID,
			FishType:          item.FishType,
			Quantity:          item.Quantity,
			Unit:              item.Unit,
			Status:            item.Status.String(),
			HighestBid:        highestBid,
			HighestBidderID:   item.HighestBidderID,
			HighestBidderName: item.HighestBidderName,
			SortOrder:         item.SortOrder,
			CreatedAt:         item.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func parseTime(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("15:04:05", *s)
	if err != nil {
		t, err = time.Parse("15:04", *s)
		if err != nil {
			return nil
		}
	}
	return new(t)
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("15:04:05")
	return new(s)
}

// RegisterRoutes registers the public auction handler routes to the given mux.
func (h *PublicAuctionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/auctions", h.List)
	mux.HandleFunc("GET /api/auctions/{id}", h.Get)
	mux.HandleFunc("GET /api/auctions/{id}/items", h.GetItems)
}
