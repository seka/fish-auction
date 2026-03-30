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
	mockReg := &mock.MockRegistry{
		ListBuyersUC: &mock.MockListBuyersUseCase{
			ExecuteFunc: func(_ context.Context) ([]model.Buyer, error) {
				return []model.Buyer{{Name: "B1"}}, nil
			},
		},
	}
	h := admin.NewBuyerHandler(mockReg)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/buyers", nil)
	w := httptest.NewRecorder()
	h.List(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAdminBuyerHandler_Delete(t *testing.T) {
	mockReg := &mock.MockRegistry{
		DeleteBuyerUC: &mock.MockDeleteBuyerUseCase{
			ExecuteFunc: func(_ context.Context, id int) error {
				if id == 999 {
					return errors.New("not found")
				}
				return nil
			},
		},
	}
	h := admin.NewBuyerHandler(mockReg)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequestWithContext(context.Background(), http.MethodDelete, "/buyers/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.Delete(w, req)
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", w.Code)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		req := httptest.NewRequestWithContext(context.Background(), http.MethodDelete, "/buyers/invalid", nil)
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()
		h.Delete(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}
