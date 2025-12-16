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
