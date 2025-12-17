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

func TestVenueHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateVenueUseCase{
			ExecuteFunc: func(ctx context.Context, venue *model.Venue) (*model.Venue, error) {
				venue.ID = 1
				return venue, nil
			},
		}
		mockReg := &mock.MockRegistry{CreateVenueUC: mockCreateUC}
		h := handler.NewVenueHandler(mockReg)

		reqBody := dto.CreateVenueRequest{Name: "Tokyo Market", Location: "Toyosu"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/venues", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Error_UseCase", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateVenueUseCase{
			ExecuteFunc: func(ctx context.Context, venue *model.Venue) (*model.Venue, error) {
				return nil, errors.New("db error")
			},
		}
		mockReg := &mock.MockRegistry{CreateVenueUC: mockCreateUC}
		h := handler.NewVenueHandler(mockReg)

		reqBody := dto.CreateVenueRequest{Name: "Tokyo Market", Location: "Toyosu"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/venues", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestVenueHandler_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockListUC := &mock.MockListVenuesUseCase{
			ExecuteFunc: func(ctx context.Context) ([]model.Venue, error) {
				return []model.Venue{{ID: 1, Name: "V1"}, {ID: 2, Name: "V2"}}, nil
			},
		}
		mockReg := &mock.MockRegistry{ListVenuesUC: mockListUC}
		h := handler.NewVenueHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/venues", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestVenueHandler_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockGetUC := &mock.MockGetVenueUseCase{
			ExecuteFunc: func(ctx context.Context, id int) (*model.Venue, error) {
				return &model.Venue{ID: 1, Name: "V1"}, nil
			},
		}
		mockReg := &mock.MockRegistry{GetVenueUC: mockGetUC}
		h := handler.NewVenueHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/venues/1", nil)
		w := httptest.NewRecorder()

		h.Get(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockGetUC := &mock.MockGetVenueUseCase{
			ExecuteFunc: func(ctx context.Context, id int) (*model.Venue, error) {
				return nil, errors.New("not found")
			},
		}
		mockReg := &mock.MockRegistry{GetVenueUC: mockGetUC}
		h := handler.NewVenueHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/venues/999", nil)
		w := httptest.NewRecorder()

		h.Get(w, req)

		if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
			t.Errorf("expected error status, got %d", w.Code)
		}
	})
}

func TestVenueHandler_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUpdateUC := &mock.MockUpdateVenueUseCase{
			ExecuteFunc: func(ctx context.Context, venue *model.Venue) error {
				if venue.ID != 1 {
					t.Errorf("expected ID 1, got %d", venue.ID)
				}
				return nil
			},
		}
		mockReg := &mock.MockRegistry{UpdateVenueUC: mockUpdateUC}
		h := handler.NewVenueHandler(mockReg)

		reqBody := dto.UpdateVenueRequest{Name: "Updated V1"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/api/venues/1", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Update(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestVenueHandler_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDeleteUC := &mock.MockDeleteVenueUseCase{
			ExecuteFunc: func(ctx context.Context, id int) error {
				if id != 1 {
					t.Errorf("expected ID 1, got %d", id)
				}
				return nil
			},
		}
		mockReg := &mock.MockRegistry{DeleteVenueUC: mockDeleteUC}
		h := handler.NewVenueHandler(mockReg)

		req := httptest.NewRequest(http.MethodDelete, "/api/venues/1", nil)
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", w.Code)
		}
	})
}

func TestVenueHandler_RegisterRoutes(t *testing.T) {
	t.Run("MethodNotAllowed", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := handler.NewVenueHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequest(http.MethodPut, "/api/venues", nil) // Create/List only on root
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}
