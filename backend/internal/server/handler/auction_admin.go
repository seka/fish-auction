package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
)

// AdminAuctionHandler handles admin HTTP requests related to auctions.
type AdminAuctionHandler struct {
	createUseCase       auction.CreateAuctionUseCase
	updateUseCase       auction.UpdateAuctionUseCase
	updateStatusUseCase auction.UpdateAuctionStatusUseCase
	deleteUseCase       auction.DeleteAuctionUseCase
}

// NewAdminAuctionHandler creates a new AdminAuctionHandler instance.
func NewAdminAuctionHandler(r registry.UseCase) *AdminAuctionHandler {
	return &AdminAuctionHandler{
		createUseCase:       r.NewCreateAuctionUseCase(),
		updateUseCase:       r.NewUpdateAuctionUseCase(),
		updateStatusUseCase: r.NewUpdateAuctionStatusUseCase(),
		deleteUseCase:       r.NewDeleteAuctionUseCase(),
	}
}

// Create handles the auction creation request.
func (h *AdminAuctionHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	auc := &model.Auction{
		VenueID: req.VenueID,
		Status:  model.AuctionStatus(req.Status),
		Period:  model.NewAuctionPeriod(auctionDate, parseTime(req.StartTime), parseTime(req.EndTime)),
	}

	if auc.Status == "" {
		auc.Status = model.AuctionStatusScheduled
	}

	created, err := h.createUseCase.Execute(r.Context(), auc)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := dto.AuctionResponse{
		ID:          created.ID,
		VenueID:     created.VenueID,
		AuctionDate: created.Period.AuctionDate.Format("2006-01-02"),
		StartTime:   formatTime(created.Period.StartAt),
		EndTime:     formatTime(created.Period.EndAt),
		Status:      string(created.Status),
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		util.HandleError(w, err)
		return
	}
}

// Update handles the request to update a specific auction.
func (h *AdminAuctionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
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

	auc := &model.Auction{
		ID:      id,
		VenueID: req.VenueID,
		Status:  model.AuctionStatus(req.Status),
		Period:  model.NewAuctionPeriod(auctionDate, parseTime(req.StartTime), parseTime(req.EndTime)),
	}

	if err := h.updateUseCase.Execute(r.Context(), auc); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateStatus handles the request to update the status of an auction.
func (h *AdminAuctionHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
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

// Delete handles the request to delete a specific auction.
func (h *AdminAuctionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
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

// RegisterRoutes registers the admin auction handler routes to the given mux.
func (h *AdminAuctionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auctions", h.Create)
	mux.HandleFunc("PUT /auctions/{id}", h.Update)
	mux.HandleFunc("PATCH /auctions/{id}/status", h.UpdateStatus)
	mux.HandleFunc("DELETE /auctions/{id}", h.Delete)
}
