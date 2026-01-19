package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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

		mockReg := &mock.MockRegistry{
			CreateItemUC:          mockCreateUC,
			ListItemsUC:           mockListUC,
			UpdateItemUC:          &mock.MockUpdateItemUseCase{},
			DeleteItemUC:          &mock.MockDeleteItemUseCase{},
			UpdateItemSortOrderUC: &mock.MockUpdateItemSortOrderUseCase{},
		}

		h := handler.NewItemHandler(mockReg)

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

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var resp dto.ItemResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.ID != 1 {
			t.Errorf("expected ID 1, got %d", resp.ID)
		}
	})

	t.Run("Error_InvalidJSON", func(t *testing.T) {
		mockReg := &mock.MockRegistry{
			CreateItemUC:          &mock.MockCreateItemUseCase{},
			ListItemsUC:           &mock.MockListItemsUseCase{},
			UpdateItemUC:          &mock.MockUpdateItemUseCase{},
			DeleteItemUC:          &mock.MockDeleteItemUseCase{},
			UpdateItemSortOrderUC: &mock.MockUpdateItemSortOrderUseCase{},
		}

		h := handler.NewItemHandler(mockReg)

		req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestItemHandler_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockListUC := &mock.MockListItemsUseCase{
			ExecuteFunc: func(ctx context.Context, status string) ([]model.AuctionItem, error) {
				return []model.AuctionItem{
					{ID: 1, FishType: "Tuna", Quantity: 10, Unit: "kg", Status: model.ItemStatusAvailable},
				}, nil
			},
		}

		mockReg := &mock.MockRegistry{
			CreateItemUC:          &mock.MockCreateItemUseCase{},
			ListItemsUC:           mockListUC,
			UpdateItemUC:          &mock.MockUpdateItemUseCase{},
			DeleteItemUC:          &mock.MockDeleteItemUseCase{},
			UpdateItemSortOrderUC: &mock.MockUpdateItemSortOrderUseCase{},
		}

		h := handler.NewItemHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/items?status=Available", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestItemHandler_RegisterRoutes(t *testing.T) {
	t.Run("MethodNotAllowed", func(t *testing.T) {
		mockReg := &mock.MockRegistry{
			CreateItemUC:          &mock.MockCreateItemUseCase{},
			ListItemsUC:           &mock.MockListItemsUseCase{},
			UpdateItemUC:          &mock.MockUpdateItemUseCase{},
			DeleteItemUC:          &mock.MockDeleteItemUseCase{},
			UpdateItemSortOrderUC: &mock.MockUpdateItemSortOrderUseCase{},
		}
		h := handler.NewItemHandler(mockReg)
		r := mux.NewRouter()
		authMiddleware := func(next http.Handler) http.Handler { return next }
		h.RegisterRoutes(r, authMiddleware)

		req := httptest.NewRequest(http.MethodPatch, "/api/items", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}
