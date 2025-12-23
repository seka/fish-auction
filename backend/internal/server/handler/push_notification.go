package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/notification"
)

type PushHandler struct {
	pushUseCase notification.PushNotificationUseCase
}

func NewPushHandler(r registry.UseCase) *PushHandler {
	return &PushHandler{
		pushUseCase: r.NewPushNotificationUseCase(),
	}
}

func (h *PushHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req dto.SubscribePushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, err)
		return
	}

	// Get buyer ID from context (authenticated user)
	buyerID, ok := r.Context().Value("buyer_id").(int)
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

func (h *PushHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/push/subscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Subscribe(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
