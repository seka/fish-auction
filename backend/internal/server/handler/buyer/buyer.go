package buyer

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer/response"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
)

// BuyerHandler handles buyer HTTP requests related to their account and purchases.
type BuyerHandler struct {
	getMeUseCase        buyer.GetBuyerUseCase
	getPurchasesUseCase buyer.GetBuyerPurchasesUseCase
	getAuctionsUseCase  buyer.GetBuyerAuctionsUseCase
	updatePassUseCase   buyer.UpdatePasswordUseCase
}

// NewBuyerHandler creates a new BuyerHandler instance.
func NewBuyerHandler(r registry.UseCase) *BuyerHandler {
	return &BuyerHandler{
		getMeUseCase:        r.NewGetBuyerUseCase(),
		getPurchasesUseCase: r.NewGetBuyerPurchasesUseCase(),
		getAuctionsUseCase:  r.NewGetBuyerAuctionsUseCase(),
		updatePassUseCase:   r.NewBuyerUpdatePasswordUseCase(),
	}
}

// GetMe handles the request to get the current authenticated buyer's info.
func (h *BuyerHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	b, err := h.getMeUseCase.Execute(r.Context(), buyerID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, response.Me{
		Authenticated: true,
		BuyerID:       b.ID,
		Name:          b.Name,
	})
}

// GetPurchases handles the request to get the buyer's purchases.
func (h *BuyerHandler) GetPurchases(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	purchases, err := h.getPurchasesUseCase.Execute(r.Context(), buyerID)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]response.Purchase, len(purchases))
	for i, p := range purchases {
		resp[i] = response.Purchase{
			ID:          p.ID,
			ItemID:      p.ItemID,
			FishType:    p.FishType,
			Quantity:    p.Quantity,
			Unit:        p.Unit,
			Price:       p.Price,
			BuyerID:     p.BuyerID,
			AuctionID:   p.AuctionID,
			AuctionDate: p.AuctionDate,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		}
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

// GetAuctions handles the request to get auctions the buyer participated in.
func (h *BuyerHandler) GetAuctions(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	auctions, err := h.getAuctionsUseCase.Execute(r.Context(), buyerID)
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
			StartTime:   formatTime(a.Period.StartAt),
			EndTime:     formatTime(a.Period.EndAt),
			Status:      string(a.Status),
			CreatedAt:   a.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   a.UpdatedAt.Format(time.RFC3339),
		}
	}

	util.WriteJSON(w, http.StatusOK, resp)
}

// UpdatePassword handles the request to update the buyer's password.
func (h *BuyerHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req request.UpdatePassword
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	if err := h.updatePassUseCase.Execute(r.Context(), buyerID, req.CurrentPassword, req.NewPassword); err != nil {
		util.HandleError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, response.Message{Message: "Password updated successfully"})
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("15:04:05")
	return &s
}

// RegisterRoutes registers the buyer account handler routes to the given mux.
func (h *BuyerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /me", h.GetMe)
	mux.HandleFunc("GET /purchases", h.GetPurchases)
	mux.HandleFunc("GET /auctions", h.GetAuctions)
	mux.HandleFunc("PUT /password", h.UpdatePassword)
}
