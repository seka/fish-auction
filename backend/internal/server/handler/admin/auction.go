package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/auction"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
)

// AuctionHandler handles admin HTTP requests related to auctions.
type AuctionHandler struct {
	createUseCase       auction.CreateAuctionUseCase
	updateUseCase       auction.UpdateAuctionUseCase
	updateStatusUseCase auction.UpdateAuctionStatusUseCase
	deleteUseCase       auction.DeleteAuctionUseCase
	reorderItemsUseCase item.ReorderItemsUseCase
}

// NewAuctionHandler creates a new AuctionHandler instance.
func NewAuctionHandler(r registry.UseCase) *AuctionHandler {
	return &AuctionHandler{
		createUseCase:       r.NewCreateAuctionUseCase(),
		updateUseCase:       r.NewUpdateAuctionUseCase(),
		updateStatusUseCase: r.NewUpdateAuctionStatusUseCase(),
		deleteUseCase:       r.NewDeleteAuctionUseCase(),
		reorderItemsUseCase: r.NewReorderItemsUseCase(),
	}
}

// Create handles the auction creation request.
func (h *AuctionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateAuction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	startAt, err := parseTimestamp(req.StartAt)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid start_at format (RFC3339)")
		return
	}
	endAt, err := parseTimestamp(req.EndAt)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid end_at format (RFC3339)")
		return
	}

	auc := &model.Auction{
		VenueID: req.VenueID,
		Status:  model.AuctionStatus(req.Status),
		Period:  model.NewAuctionPeriod(startAt, endAt),
	}

	if auc.Status == "" {
		auc.Status = model.AuctionStatusScheduled
	}

	created, err := h.createUseCase.Execute(r.Context(), auc)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, h.toResponse(created))
}

// Update handles the request to update a specific auction.
func (h *AuctionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req request.UpdateAuction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	startAt, err := parseTimestamp(req.StartAt)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid start_at format (RFC3339)")
		return
	}
	endAt, err := parseTimestamp(req.EndAt)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid end_at format (RFC3339)")
		return
	}

	auc := &model.Auction{
		ID:      id,
		VenueID: req.VenueID,
		Status:  model.AuctionStatus(req.Status),
		Period:  model.NewAuctionPeriod(startAt, endAt),
	}

	if err := h.updateUseCase.Execute(r.Context(), auc); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateStatus handles the request to update the status of an auction.
func (h *AuctionHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req request.UpdateAuctionStatus
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
func (h *AuctionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Reorder handles the request to reorder items within an auction.
func (h *AuctionHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	auctionID, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid auction ID")
		return
	}
	var req request.ReorderItems
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := h.reorderItemsUseCase.Execute(r.Context(), auctionID, req.IDs); err != nil {
		util.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuctionHandler) toResponse(a *model.Auction) response.Auction {
	return response.Auction{
		ID:        a.ID,
		VenueID:   a.VenueID,
		StartAt:   util.FormatTimestamp(a.Period.StartAt),
		EndAt:     util.FormatTimestamp(a.Period.EndAt),
		Status:    string(a.Status),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func parseTimestamp(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// RegisterRoutes registers the admin auction handler routes to the given mux.
func (h *AuctionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auctions", h.Create)
	mux.HandleFunc("PUT /auctions/{id}", h.Update)
	mux.HandleFunc("PATCH /auctions/{id}/status", h.UpdateStatus)
	mux.HandleFunc("DELETE /auctions/{id}", h.Delete)
	mux.HandleFunc("PUT /auctions/{id}/reorder", h.Reorder)
}
