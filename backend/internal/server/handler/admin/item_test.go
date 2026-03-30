package admin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	admin "github.com/seka/fish-auction/backend/internal/server/handler/admin"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/request"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/response"
	public "github.com/seka/fish-auction/backend/internal/server/handler/public"
	publicResponse "github.com/seka/fish-auction/backend/internal/server/handler/public/response"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAdminItemHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateItemUseCase{
			ExecuteFunc: func(_ context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
				item.ID = 1
				item.Status = model.ItemStatusAvailable
				return item, nil
			},
		}
		mockReg := &mock.MockRegistry{
			CreateItemUC: mockCreateUC,
		}

		h := admin.NewItemHandler(mockReg)

		reqBody := request.CreateItem{
			FishermanID: 1,
			FishType:    "Tuna",
			Quantity:    10,
			Unit:        "kg",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/items", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var resp response.Item
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.ID != 1 {
			t.Errorf("expected ID 1, got %d", resp.ID)
		}
	})

	t.Run("Error_InvalidJSON", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := admin.NewItemHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/items", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestPublicItemHandler_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockListUC := &mock.MockListItemsUseCase{
			ExecuteFunc: func(_ context.Context, _ string) ([]model.AuctionItem, error) {
				return []model.AuctionItem{
					{ID: 1, FishType: "Tuna", Quantity: 10, Unit: "kg", Status: model.ItemStatusAvailable},
				}, nil
			},
		}

		mockReg := &mock.MockRegistry{
			ListItemsUC: mockListUC,
		}

		h := public.NewItemHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/items?status=Available", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp []publicResponse.Item
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(resp) != 1 {
			t.Errorf("expected 1 item, got %d", len(resp))
		}
	})
}

func TestItemHandler_RegisterRoutes(t *testing.T) {
	mockReg := &mock.MockRegistry{
		ListItemsUC:  &mock.MockListItemsUseCase{ExecuteFunc: func(_ context.Context, _ string) ([]model.AuctionItem, error) { return []model.AuctionItem{}, nil }},
		CreateItemUC: &mock.MockCreateItemUseCase{ExecuteFunc: func(_ context.Context, i *model.AuctionItem) (*model.AuctionItem, error) { return i, nil }},
	}

	t.Run("Public", func(t *testing.T) {
		h := public.NewItemHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/items", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Admin", func(t *testing.T) {
		h := admin.NewItemHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/items", bytes.NewReader([]byte("{}")))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})
}
