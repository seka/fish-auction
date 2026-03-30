package buyer

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer/response"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/bid"
)

// BidHandler handles buyer HTTP requests related to bidding.
type BidHandler struct {
	createUseCase bid.CreateBidUseCase
}

// NewBidHandler creates a new BidHandler instance.
func NewBidHandler(r registry.UseCase) *BidHandler {
	return &BidHandler{
		createUseCase: r.NewCreateBidUseCase(),
	}
}

// Create handles the bid creation request.
func (h *BidHandler) Create(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req request.CreateBid
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	b := &model.Bid{
		ItemID:  req.ItemID,
		BuyerID: buyerID,
		Price:   model.NewBidPrice(req.Price),
	}

	created, err := h.createUseCase.Execute(r.Context(), b)
	if err != nil {
		util.HandleError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, response.Bid{
		ID:        created.ID,
		ItemID:    created.ItemID,
		BuyerID:   created.BuyerID,
		Price:     created.Price.Amount(),
		CreatedAt: created.CreatedAt,
	})
}

// RegisterRoutes registers the buyer bid handler routes to the given mux.
func (h *BidHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /bids", h.Create)
}
