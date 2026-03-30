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
	"github.com/seka/fish-auction/backend/internal/server/handler/admin"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/request"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAdminBuyerHandler_Create(t *testing.T) {
	type testCase struct {
		name       string
		body       any
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			body: request.CreateBuyer{
				Name:         "Buyer 1",
				Email:        "buyer@example.com",
				Password:     "password",
				Organization: "Org 1",
				ContactInfo:  "Contact 1",
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.CreateBuyerUC = &mock.MockCreateBuyerUseCase{
					ExecuteFunc: func(_ context.Context, name, _, _, organization, _ string) (*model.Buyer, error) {
						return &model.Buyer{ID: 1, Name: name, Organization: organization}, nil
					},
				}
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "InvalidJSON",
			body:       "invalid-json",
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "UseCaseError",
			body: request.CreateBuyer{Email: "buyer@example.com"},
			mockSetup: func(r *mock.MockRegistry) {
				r.CreateBuyerUC = &mock.MockCreateBuyerUseCase{
					ExecuteFunc: func(_ context.Context, _, _, _, _, _ string) (*model.Buyer, error) {
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
			h := admin.NewBuyerHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/buyers", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAdminBuyerHandler_List(t *testing.T) {
	type testCase struct {
		name       string
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			mockSetup: func(r *mock.MockRegistry) {
				r.ListBuyersUC = &mock.MockListBuyersUseCase{
					ExecuteFunc: func(_ context.Context) ([]model.Buyer, error) {
						return []model.Buyer{{Name: "B1"}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := admin.NewBuyerHandler(mockReg)
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/buyers", nil)
			w := httptest.NewRecorder()
			h.List(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}
