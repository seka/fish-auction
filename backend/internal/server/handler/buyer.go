package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

// BuyerHandler handles HTTP requests related to buyers.
type BuyerHandler struct {
	loginUseCase          buyer.LoginBuyerUseCase
	getPurchasesUseCase   buyer.GetBuyerPurchasesUseCase
	getAuctionsUseCase    buyer.GetBuyerAuctionsUseCase
	updatePasswordUseCase buyer.UpdatePasswordUseCase
	getBuyerUseCase       buyer.GetBuyerUseCase
	sessionRepo           repository.SessionRepository
}

// NewBuyerHandler creates a new BuyerHandler instance.
func NewBuyerHandler(r registry.UseCase, sessionRepo repository.SessionRepository) *BuyerHandler {
	return &BuyerHandler{
		loginUseCase:          r.NewLoginBuyerUseCase(),
		getPurchasesUseCase:   r.NewGetBuyerPurchasesUseCase(),
		getAuctionsUseCase:    r.NewGetBuyerAuctionsUseCase(),
		updatePasswordUseCase: r.NewBuyerUpdatePasswordUseCase(),
		getBuyerUseCase:       r.NewGetBuyerUseCase(),
		sessionRepo:           sessionRepo,
	}
}

// Login handles the buyer login request.
func (h *BuyerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginBuyerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	buy, err := h.loginUseCase.Execute(r.Context(), req.Email, req.Password)
	if err != nil {
		util.HandleError(w, err)
		return
	}
	if buy == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID, err := h.sessionRepo.Create(r.Context(), buy.ID, model.SessionRoleBuyer)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	setSessionCookie(w, "buyer_session", sessionID)

	resp := dto.BuyerResponse{ID: buy.ID, Name: buy.Name}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Logout handles the buyer logout request.
func (h *BuyerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("buyer_session"); err == nil {
		if err := h.sessionRepo.Delete(r.Context(), cookie.Value); err != nil {
			util.HandleError(w, err)
			return
		}
	}

	clearSessionCookie(w, "buyer_session")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dto.MessageResponse{Message: "Logged out"}); err != nil {
		util.HandleError(w, err)
		return
	}
}

// GetCurrentBuyer handles the request to get the currently authenticated buyer.
func (h *BuyerHandler) GetCurrentBuyer(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	buy, err := h.getBuyerUseCase.Execute(r.Context(), buyerID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dto.BuyerMeResponse{
		Authenticated: true,
		BuyerID:       buyerID,
		Name:          buy.Name,
	}); err != nil {
		util.HandleError(w, err)
		return
	}
}

// GetMyPurchases handles the request to get the purchases of the authenticated buyer.
func (h *BuyerHandler) GetMyPurchases(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	purchases, err := h.getPurchasesUseCase.Execute(r.Context(), buyerID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	// Convert to DTOs
	resp := make([]dto.PurchaseResponse, len(purchases))
	for i, p := range purchases {
		resp[i] = dto.PurchaseResponse{
			ID:          p.ID,
			ItemID:      p.ItemID,
			FishType:    p.FishType,
			Quantity:    p.Quantity,
			Unit:        p.Unit,
			Price:       p.Price,
			BuyerID:     p.BuyerID,
			AuctionID:   p.AuctionID,
			AuctionDate: p.AuctionDate,
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// GetMyAuctions handles the request to get the auctions of the authenticated buyer.
func (h *BuyerHandler) GetMyAuctions(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	auctions, err := h.getAuctionsUseCase.Execute(r.Context(), buyerID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	// Convert to DTOs
	resp := make([]dto.AuctionResponse, len(auctions))
	for i, a := range auctions {
		var startTime, endTime *string
		if a.Period.StartAt != nil {
			s := a.Period.StartAt.Format("15:04:05")
			startTime = &s
		}
		if a.Period.EndAt != nil {
			e := a.Period.EndAt.Format("15:04:05")
			endTime = &e
		}

		resp[i] = dto.AuctionResponse{
			ID:          a.ID,
			VenueID:     a.VenueID,
			AuctionDate: a.Period.AuctionDate.Format("2006-01-02"),
			StartTime:   startTime,
			EndTime:     endTime,
			Status:      string(a.Status),
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// UpdatePassword handles the password update request for a buyer.
func (h *BuyerHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	if err := h.updatePasswordUseCase.Execute(r.Context(), buyerID, req.CurrentPassword, req.NewPassword); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dto.MessageResponse{Message: "Password updated successfully"}); err != nil {
		util.HandleError(w, err)
		return
	}
}

// RegisterPublicRoutes registers the public buyer handler routes to the given mux.
func (h *BuyerHandler) RegisterPublicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/buyer/login", h.Login)
	mux.HandleFunc("POST /api/buyer/logout", h.Logout)
}

// RegisterBuyerRoutes registers the authenticated buyer handler routes to the given mux.
func (h *BuyerHandler) RegisterBuyerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /me", h.GetCurrentBuyer)
	mux.HandleFunc("GET /me/purchases", h.GetMyPurchases)
	mux.HandleFunc("GET /me/auctions", h.GetMyAuctions)
	mux.HandleFunc("PUT /password", h.UpdatePassword)
}
