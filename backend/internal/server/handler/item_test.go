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

func TestItemHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateItemUseCase{
			ExecuteFunc: func(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
				item.ID = 1
				item.Status = model.ItemStatusAvailable
				return item, nil
			},
		}
		mockListUC := &mock.MockListItemsUseCase{}

		h := handler.NewItemHandler(mockCreateUC, mockListUC)

		reqBody := dto.CreateItemRequest{
			FishermanID: 1,
			FishType:    "Tuna",
			Quantity:    10,
			Unit:        "kg",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp dto.ItemResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.ID != 1 {
			t.Errorf("expected ID 1, got %d", resp.ID)
		}
		if resp.FishType != "Tuna" {
			t.Errorf("expected FishType Tuna, got %s", resp.FishType)
		}
	})

	t.Run("Error_InvalidJSON", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateItemUseCase{}
		mockListUC := &mock.MockListItemsUseCase{}

		h := handler.NewItemHandler(mockCreateUC, mockListUC)

		req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})

	t.Run("Error_UseCaseFails", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateItemUseCase{
			ExecuteFunc: func(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
				return nil, errors.New("db error")
			},
		}
		mockListUC := &mock.MockListItemsUseCase{}

		h := handler.NewItemHandler(mockCreateUC, mockListUC)

		reqBody := dto.CreateItemRequest{
			FishermanID: 1,
			FishType:    "Tuna",
			Quantity:    10,
			Unit:        "kg",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestItemHandler_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateItemUseCase{}
		mockListUC := &mock.MockListItemsUseCase{
			ExecuteFunc: func(ctx context.Context, status string) ([]model.AuctionItem, error) {
				return []model.AuctionItem{
					{ID: 1, FishType: "Tuna", Quantity: 10, Unit: "kg", Status: model.ItemStatusAvailable},
					{ID: 2, FishType: "Salmon", Quantity: 5, Unit: "kg", Status: model.ItemStatusAvailable},
				}, nil
			},
		}

		h := handler.NewItemHandler(mockCreateUC, mockListUC)

		req := httptest.NewRequest(http.MethodGet, "/api/items?status=Available", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp []dto.ItemResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(resp) != 2 {
			t.Errorf("expected 2 items, got %d", len(resp))
		}
	})
}
