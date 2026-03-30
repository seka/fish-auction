package admin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	admin "github.com/seka/fish-auction/backend/internal/server/handler/admin"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/request"
	public "github.com/seka/fish-auction/backend/internal/server/handler/public"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAdminVenueHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateVenueUseCase{
			ExecuteFunc: func(_ context.Context, venue *model.Venue) (*model.Venue, error) {
				venue.ID = 1
				return venue, nil
			},
		}
		mockReg := &mock.MockRegistry{CreateVenueUC: mockCreateUC}
		h := admin.NewVenueHandler(mockReg)

		reqBody := request.CreateVenue{Name: "Tokyo Market", Location: "Toyosu"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/venues", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("Error_UseCase", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateVenueUseCase{
			ExecuteFunc: func(_ context.Context, _ *model.Venue) (*model.Venue, error) {
				return nil, errors.New("db error")
			},
		}
		mockReg := &mock.MockRegistry{CreateVenueUC: mockCreateUC}
		h := admin.NewVenueHandler(mockReg)

		reqBody := request.CreateVenue{Name: "Tokyo Market", Location: "Toyosu"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/venues", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestPublicVenueHandler_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockListUC := &mock.MockListVenuesUseCase{
			ExecuteFunc: func(_ context.Context) ([]model.Venue, error) {
				return []model.Venue{{ID: 1, Name: "V1"}, {ID: 2, Name: "V2"}}, nil
			},
		}
		mockReg := &mock.MockRegistry{ListVenuesUC: mockListUC}
		h := public.NewVenueHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/venues", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestPublicVenueHandler_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockGetUC := &mock.MockGetVenueUseCase{
			ExecuteFunc: func(_ context.Context, _ int) (*model.Venue, error) {
				return &model.Venue{ID: 1, Name: "V1"}, nil
			},
		}
		mockReg := &mock.MockRegistry{GetVenueUC: mockGetUC}
		h := public.NewVenueHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/venues/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		h.Get(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockGetUC := &mock.MockGetVenueUseCase{
			ExecuteFunc: func(_ context.Context, _ int) (*model.Venue, error) {
				return nil, errors.New("not found")
			},
		}
		mockReg := &mock.MockRegistry{GetVenueUC: mockGetUC}
		h := public.NewVenueHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/venues/999", nil)
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()

		h.Get(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected error status 500, got %d", w.Code)
		}
	})
}

func TestAdminVenueHandler_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUpdateUC := &mock.MockUpdateVenueUseCase{
			ExecuteFunc: func(_ context.Context, venue *model.Venue) error {
				if venue.ID != 1 {
					t.Errorf("expected ID 1, got %d", venue.ID)
				}
				return nil
			},
		}
		mockReg := &mock.MockRegistry{UpdateVenueUC: mockUpdateUC}
		h := admin.NewVenueHandler(mockReg)

		reqBody := request.UpdateVenue{Name: "Updated V1"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/venues/1", bytes.NewReader(body))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		h.Update(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestAdminVenueHandler_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDeleteUC := &mock.MockDeleteVenueUseCase{
			ExecuteFunc: func(_ context.Context, id int) error {
				if id != 1 {
					t.Errorf("expected ID 1, got %d", id)
				}
				return nil
			},
		}
		mockReg := &mock.MockRegistry{DeleteVenueUC: mockDeleteUC}
		h := admin.NewVenueHandler(mockReg)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodDelete, "/venues/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		h.Delete(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", w.Code)
		}
	})
}

func TestVenueHandler_RegisterRoutes(t *testing.T) {
	mockReg := &mock.MockRegistry{
		ListVenuesUC:  &mock.MockListVenuesUseCase{ExecuteFunc: func(_ context.Context) ([]model.Venue, error) { return []model.Venue{}, nil }},
		CreateVenueUC: &mock.MockCreateVenueUseCase{ExecuteFunc: func(_ context.Context, v *model.Venue) (*model.Venue, error) { v.ID = 1; return v, nil }},
		GetVenueUC:    &mock.MockGetVenueUseCase{ExecuteFunc: func(_ context.Context, _ int) (*model.Venue, error) { return &model.Venue{ID: 1}, nil }},
	}

	t.Run("Public", func(t *testing.T) {
		h := public.NewVenueHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/venues", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Admin", func(t *testing.T) {
		h := admin.NewVenueHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/venues", bytes.NewReader([]byte("{}")))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", w.Code)
		}
	})
}
