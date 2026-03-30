package public_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/handler/public"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestPublicVenueHandler_List(t *testing.T) {
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
}

func TestPublicVenueHandler_Get(t *testing.T) {
	type testCase struct {
		name       string
		idStr      string
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name:  "Success",
			idStr: "1",
			mockSetup: func(r *mock.MockRegistry) {
				r.GetVenueUC = &mock.MockGetVenueUseCase{
					ExecuteFunc: func(_ context.Context, _ int) (*model.Venue, error) {
						return &model.Venue{ID: 1, Name: "V1"}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidID",
			idStr:      "invalid",
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "NotFound",
			idStr: "999",
			mockSetup: func(r *mock.MockRegistry) {
				r.GetVenueUC = &mock.MockGetVenueUseCase{
					ExecuteFunc: func(_ context.Context, _ int) (*model.Venue, error) {
						return nil, errors.New("not found")
					},
				}
			},
			wantStatus: http.StatusInternalServerError, // HandleError maps generic errors to 500
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := public.NewVenueHandler(mockReg)

			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/venues/"+tc.idStr, nil)
			req.SetPathValue("id", tc.idStr)
			w := httptest.NewRecorder()

			h.Get(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}
