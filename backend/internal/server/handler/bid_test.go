package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestBidHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateBidUseCase{
			ExecuteFunc: func(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
				if bid.BuyerID != 1 {
					t.Errorf("expected buyerID 1, got %d", bid.BuyerID)
				}
				bid.ID = 1
				return bid, nil
			},
		}
		mockReg := &mock.MockRegistry{CreateBidUC: mockCreateUC}
		h := handler.NewBidHandler(mockReg)

		reqBody := dto.CreateBidRequest{ItemID: 10, Price: 500}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/bids", bytes.NewReader(body))

		// Inject buyer_id into context (simulating middleware)
		ctx := context.WithValue(req.Context(), "buyer_id", 1)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("Unauthorized_NoContext", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := handler.NewBidHandler(mockReg)

		reqBody := dto.CreateBidRequest{ItemID: 10, Price: 500}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/bids", bytes.NewReader(body))
		// No Context

		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateBidUseCase{
			ExecuteFunc: func(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
				return nil, errors.New("validation failed")
			},
		}
		mockReg := &mock.MockRegistry{CreateBidUC: mockCreateUC}
		h := handler.NewBidHandler(mockReg)

		reqBody := dto.CreateBidRequest{ItemID: 10, Price: 500}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/bids", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), "buyer_id", 1)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}
