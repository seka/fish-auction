package admin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/request"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestFishermanHandler(t *testing.T) {
	mockReg := &mock.MockRegistry{
		CreateFishermanUC: &mock.MockCreateFishermanUseCase{
			ExecuteFunc: func(_ context.Context, n string) (*model.Fisherman, error) {
				return &model.Fisherman{ID: 1, Name: n}, nil
			},
		},
		ListFishermenUC: &mock.MockListFishermenUseCase{
			ExecuteFunc: func(_ context.Context) ([]model.Fisherman, error) {
				return []model.Fisherman{{ID: 1, Name: "F1"}}, nil
			},
		},
	}
	h := admin.NewFishermanHandler(mockReg)

	t.Run("Create_Success", func(t *testing.T) {
		reqBody := request.CreateFisherman{Name: "F1"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/fishermen", bytes.NewReader(body))
		w := httptest.NewRecorder()
		h.Create(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})

	t.Run("List_Success", func(t *testing.T) {
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/fishermen", nil)
		w := httptest.NewRecorder()
		h.List(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("RegisterRoutes", func(t *testing.T) {
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/fishermen", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}
