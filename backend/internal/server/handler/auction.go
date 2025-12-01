package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

type AuctionHandler struct {
	createUseCase       auction.CreateAuctionUseCase
	listUseCase         auction.ListAuctionsUseCase
	getUseCase          auction.GetAuctionUseCase
	getItemsUseCase     auction.GetAuctionItemsUseCase
	updateUseCase       auction.UpdateAuctionUseCase
	updateStatusUseCase auction.UpdateAuctionStatusUseCase
	deleteUseCase       auction.DeleteAuctionUseCase
}

func NewAuctionHandler(r registry.UseCase) *AuctionHandler {
	return &AuctionHandler{
		createUseCase:       r.NewCreateAuctionUseCase(),
		listUseCase:         r.NewListAuctionsUseCase(),
		getUseCase:          r.NewGetAuctionUseCase(),
		getItemsUseCase:     r.NewGetAuctionItemsUseCase(),
		updateUseCase:       r.NewUpdateAuctionUseCase(),
		updateStatusUseCase: r.NewUpdateAuctionStatusUseCase(),
		deleteUseCase:       r.NewDeleteAuctionUseCase(),
	}
}

func parseTime(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("15:04:05", *s)
	if err != nil {
		// Try parsing without seconds
		t, err = time.Parse("15:04", *s)
		if err != nil {
			return nil
		}
	}
	return &t
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("15:04:05")
	return &s
}

func (h *AuctionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAuctionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	auctionDate, err := time.Parse("2006-01-02", req.AuctionDate)
	if err != nil {
		http.Error(w, "Invalid date format (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	auction := &model.Auction{
		VenueID:     req.VenueID,
		AuctionDate: auctionDate,
		StartTime:   parseTime(req.StartTime),
		EndTime:     parseTime(req.EndTime),
		Status:      model.AuctionStatus(req.Status),
	}

	if auction.Status == "" {
		auction.Status = model.AuctionStatusScheduled
	}

	created, err := h.createUseCase.Execute(r.Context(), auction)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := dto.AuctionResponse{
		ID:          created.ID,
		VenueID:     created.VenueID,
		AuctionDate: created.AuctionDate.Format("2006-01-02"),
		StartTime:   formatTime(created.StartTime),
		EndTime:     formatTime(created.EndTime),
		Status:      string(created.Status),
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

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

	resp := make([]dto.AuctionResponse, len(auctions))
	for i, a := range auctions {
		resp[i] = dto.AuctionResponse{
			ID:          a.ID,
			VenueID:     a.VenueID,
			AuctionDate: a.AuctionDate.Format("2006-01-02"),
			StartTime:   formatTime(a.StartTime),
			EndTime:     formatTime(a.EndTime),
			Status:      string(a.Status),
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuctionHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/auctions/")
	// Handle /items suffix
	if strings.HasSuffix(idStr, "/items") {
		h.GetItems(w, r)
		return
	}
	// Handle /status suffix
	if strings.HasSuffix(idStr, "/status") {
		// This should be handled by PATCH method in RegisterRoutes, but just in case
		return
	}

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
		AuctionDate: a.AuctionDate.Format("2006-01-02"),
		StartTime:   formatTime(a.StartTime),
		EndTime:     formatTime(a.EndTime),
		Status:      string(a.Status),
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuctionHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/auctions/"), "/items")
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
	for i, item := range items {
		resp[i] = dto.ItemResponse{
			ID:          item.ID,
			AuctionID:   item.AuctionID,
			FishermanID: item.FishermanID,
			FishType:    item.FishType,
			Quantity:    item.Quantity,
			Unit:        item.Unit,
			Status:      item.Status.String(),
			HighestBid:  item.HighestBid,
			CreatedAt:   item.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuctionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/auctions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateAuctionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	auctionDate, err := time.Parse("2006-01-02", req.AuctionDate)
	if err != nil {
		http.Error(w, "Invalid date format (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	auction := &model.Auction{
		ID:          id,
		VenueID:     req.VenueID,
		AuctionDate: auctionDate,
		StartTime:   parseTime(req.StartTime),
		EndTime:     parseTime(req.EndTime),
		Status:      model.AuctionStatus(req.Status),
	}

	if err := h.updateUseCase.Execute(r.Context(), auction); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuctionHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/auctions/"), "/status")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateAuctionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	status := model.AuctionStatus(req.Status)
	if err := h.updateStatusUseCase.Execute(r.Context(), id, status); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuctionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/auctions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuctionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auctions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			h.List(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/auctions/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/items") {
			if r.Method == http.MethodGet {
				h.GetItems(w, r)
				return
			}
		}
		if strings.HasSuffix(r.URL.Path, "/status") {
			if r.Method == http.MethodPatch {
				h.UpdateStatus(w, r)
				return
			}
		}

		switch r.Method {
		case http.MethodGet:
			h.Get(w, r)
		case http.MethodPut:
			h.Update(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
