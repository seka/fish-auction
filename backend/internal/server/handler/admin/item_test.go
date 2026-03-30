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
		mockReg := &mock.MockRegistry{CreateItemUC: mockCreateUC}
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

func TestAdminItemHandler_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUpdateUC := &mock.MockUpdateItemUseCase{
			ExecuteFunc: func(_ context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
				return item, nil
			},
		}
		mockReg := &mock.MockRegistry{UpdateItemUC: mockUpdateUC}
		h := admin.NewItemHandler(mockReg)

		reqBody := request.UpdateItem{FishType: "Mackerel"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/items/1", bytes.NewReader(body))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		h.Update(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := admin.NewItemHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/items/invalid", bytes.NewReader([]byte("{}")))
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()

		h.Update(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestAdminItemHandler_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDeleteUC := &mock.MockDeleteItemUseCase{
			ExecuteFunc: func(_ context.Context, id int) error {
				return nil
			},
		}
		mockReg := &mock.MockRegistry{DeleteItemUC: mockDeleteUC}
		h := admin.NewItemHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodDelete, "/items/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", w.Code)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := admin.NewItemHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodDelete, "/items/invalid", nil)
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestAdminItemHandler_UpdateSortOrder(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUpdateSortUC := &mock.MockUpdateItemSortOrderUseCase{
			ExecuteFunc: func(_ context.Context, id int, sortOrder int) error {
				return nil
			},
		}
		mockReg := &mock.MockRegistry{UpdateItemSortOrderUC: mockUpdateSortUC}
		h := admin.NewItemHandler(mockReg)

		reqBody := request.UpdateItemSortOrder{SortOrder: 5}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/items/1/sort-order", bytes.NewReader(body))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		h.UpdateSortOrder(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", w.Code)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := admin.NewItemHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/items/invalid/sort-order", bytes.NewReader([]byte(`{"sort_order":1}`)))
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()

		h.UpdateSortOrder(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestItemHandler_RegisterRoutes(t *testing.T) {
	mockReg := &mock.MockRegistry{
		ListItemsUC:  &mock.MockListItemsUseCase{ExecuteFunc: func(_ context.Context, _ string) ([]model.AuctionItem, error) { return []model.AuctionItem{}, nil }},
		CreateItemUC: &mock.MockCreateItemUseCase{ExecuteFunc: func(_ context.Context, i *model.AuctionItem) (*model.AuctionItem, error) { return i, nil }},
	}

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
