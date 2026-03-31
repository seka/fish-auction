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
	"github.com/seka/fish-auction/backend/internal/usecase/notification"
)

// PushHandler handles buyer HTTP requests related to push notifications.
type PushHandler struct {
	subscribeUseCase notification.SubscribeNotificationUseCase
}

// NewPushHandler creates a new PushHandler instance.
func NewPushHandler(r registry.UseCase) *PushHandler {
	return &PushHandler{
		subscribeUseCase: r.NewSubscribeNotificationUseCase(),
	}
}

// Subscribe handles the push notification subscription request.
func (h *PushHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		util.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req request.SubscribePush
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	sub := &model.PushSubscription{
		BuyerID:  buyerID,
		Endpoint: req.Endpoint,
		P256dh:   req.Keys.P256dh,
		Auth:     req.Keys.Auth,
	}

	if err := h.subscribeUseCase.Execute(r.Context(), buyerID, sub); err != nil {
		util.HandleError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, response.Message{Message: "Subscribed successfully"})
}

// RegisterRoutes registers the buyer push notification handler routes to the given mux.
func (h *PushHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /push/subscribe", h.Subscribe)
}
