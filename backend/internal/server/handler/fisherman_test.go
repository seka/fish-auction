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

func TestFishermanHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateFishermanUseCase{
			ExecuteFunc: func(ctx context.Context, name string) (*model.Fisherman, error) {
				return &model.Fisherman{ID: 1, Name: name}, nil
			},
		}
		mockReg := &mock.MockRegistry{CreateFishermanUC: mockCreateUC}
		h := handler.NewFishermanHandler(mockReg)

		reqBody := dto.CreateFishermanRequest{Name: "Tuna Corp"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/fishermen", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Error_UseCase", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateFishermanUseCase{
			ExecuteFunc: func(ctx context.Context, name string) (*model.Fisherman, error) {
				return nil, errors.New("db error")
			},
		}
		mockReg := &mock.MockRegistry{CreateFishermanUC: mockCreateUC}
		h := handler.NewFishermanHandler(mockReg)

		reqBody := dto.CreateFishermanRequest{Name: "Tuna Corp"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/fishermen", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestFishermanHandler_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockListUC := &mock.MockListFishermenUseCase{
			ExecuteFunc: func(ctx context.Context) ([]model.Fisherman, error) {
				return []model.Fisherman{{ID: 1, Name: "F1"}, {ID: 2, Name: "F2"}}, nil
			},
		}
		mockReg := &mock.MockRegistry{ListFishermenUC: mockListUC}
		h := handler.NewFishermanHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/fishermen", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestFishermanHandler_RegisterRoutes(t *testing.T) {
	t.Run("MethodNotAllowed_PUT", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := handler.NewFishermanHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequest(http.MethodPut, "/api/fishermen", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})

	t.Run("MethodNotAllowed_POST", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := handler.NewFishermanHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequest(http.MethodPost, "/api/fishermen", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}
