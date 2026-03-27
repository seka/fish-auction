package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/notification"
)

// PushHandler handles HTTP requests related to push notifications.
type PushHandler struct {
	pushUseCase notification.PushNotificationUseCase
}

// NewPushHandler creates a new PushHandler instance.
func NewPushHandler(r registry.UseCase) *PushHandler {
	return &PushHandler{
		pushUseCase: r.NewPushNotificationUseCase(),
	}
}

// Subscribe handles the push notification subscription request.
func (h *PushHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req dto.SubscribePushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	// Get buyer ID from context (authenticated user)
	buyerID, ok := middleware.BuyerIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sub := &model.PushSubscription{
		Endpoint: req.Endpoint,
		P256dh:   req.Keys.P256dh,
		Auth:     req.Keys.Auth,
	}

	if err := h.pushUseCase.Subscribe(r.Context(), buyerID, sub); err != nil {
		util.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// RegisterRoutes registers the push notification handler routes to the given mux.
func (h *PushHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/push/subscribe", h.Subscribe)
}
