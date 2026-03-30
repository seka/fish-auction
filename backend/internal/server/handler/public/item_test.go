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

func TestPublicItemHandler_List(t *testing.T) {
	type testCase struct {
		name       string
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			mockSetup: func(r *mock.MockRegistry) {
				r.ListItemsUC = &mock.MockListItemsUseCase{
					ExecuteFunc: func(_ context.Context, _ string) ([]model.AuctionItem, error) {
						return []model.AuctionItem{{ID: 1, FishType: "Tuna"}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "UseCaseError",
			mockSetup: func(r *mock.MockRegistry) {
				r.ListItemsUC = &mock.MockListItemsUseCase{
					ExecuteFunc: func(_ context.Context, _ string) ([]model.AuctionItem, error) {
						return nil, errors.New("db error")
					},
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := public.NewItemHandler(mockReg)

			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/items", nil)
			w := httptest.NewRecorder()

			h.List(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}
