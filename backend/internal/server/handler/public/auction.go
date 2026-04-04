package public

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

// AuctionHandler handles public HTTP requests related to auctions.
type AuctionHandler struct {
	listUseCase     auction.ListAuctionsUseCase
	getUseCase      auction.GetAuctionUseCase
	getItemsUseCase auction.GetAuctionItemsUseCase
}

// NewAuctionHandler creates a new AuctionHandler instance.
func NewAuctionHandler(r registry.UseCase) *AuctionHandler {
	return &AuctionHandler{
		listUseCase:     r.NewListAuctionsUseCase(),
		getUseCase:      r.NewGetAuctionUseCase(),
		getItemsUseCase: r.NewGetAuctionItemsUseCase(),
	}
}

// List handles the request to list auctions.
func (h *AuctionHandler) List(w http.ResponseWriter, r *http.Request) {
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

	resp := make([]response.Auction, len(auctions))
	for i, a := range auctions {
		resp[i] = response.Auction{
			ID:          a.ID,
			VenueID:     a.VenueID,
			AuctionDate: a.Period.AuctionDate.Format("2006-01-02"),
			StartTime:   util.FormatTime(a.Period.StartAt),
			EndTime:     util.FormatTime(a.Period.EndAt),
			Status:      string(a.Status),
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Get handles the request to get a specific auction.
func (h *AuctionHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	a, err := h.getUseCase.Execute(r.Context(), id)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := response.Auction{
		ID:          a.ID,
		VenueID:     a.VenueID,
		AuctionDate: a.Period.AuctionDate.Format("2006-01-02"),
		StartTime:   util.FormatTime(a.Period.StartAt),
		EndTime:     util.FormatTime(a.Period.EndAt),
		Status:      string(a.Status),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// GetItems handles the request to get items for a specific auction.
func (h *AuctionHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	items, err := h.getItemsUseCase.Execute(r.Context(), id)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]response.Item, len(items))
	for i := range items {
		item := items[i]
		var highestBid *int
		if item.HighestBid != nil {
			amt := item.HighestBid.Amount()
			highestBid = &amt
		}
		resp[i] = response.Item{
			ID:          item.ID,
			AuctionID:   item.AuctionID,
			FishermanID: item.FishermanID,
			FishType:    item.FishType,
			Quantity:    item.Quantity,
			Unit:        item.Unit,
			Status:      item.Status.String(),
			HighestBid:  highestBid,
			SortOrder:   item.SortOrder,
			CreatedAt:   item.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// RegisterRoutes registers the public auction handler routes to the given mux.
func (h *AuctionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auctions", h.List)
	mux.HandleFunc("GET /auctions/{id}", h.Get)
	mux.HandleFunc("GET /auctions/{id}/items", h.GetItems)
}
