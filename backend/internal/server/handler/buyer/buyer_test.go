package buyer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer/request"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func withBuyerID(req *http.Request, buyerID int) *http.Request {
	return req.WithContext(middleware.WithBuyerID(req.Context(), buyerID))
}

func TestBuyerHandler_GetMe(t *testing.T) {
	type testCase struct {
		name        string
		withContext bool
		mockSetup   func(*mock.MockRegistry)
		wantStatus  int
	}
	tests := []testCase{
		{
			name:        "Success",
			withContext: true,
			mockSetup: func(r *mock.MockRegistry) {
				r.GetBuyerUC = &mock.MockGetBuyerUseCase{
					ExecuteFunc: func(_ context.Context, _ int) (*model.Buyer, error) {
						return &model.Buyer{ID: 1, Name: "Buyer 1"}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "NoContext_NotAuthenticated",
			withContext: false,
			mockSetup:   func(_ *mock.MockRegistry) {},
			wantStatus:  http.StatusUnauthorized,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			if tc.mockSetup != nil {
				tc.mockSetup(mockReg)
			}
			h := buyer.NewBuyerHandler(mockReg)
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/me", nil)
			if tc.withContext {
				req = withBuyerID(req, 1)
			}

			w := httptest.NewRecorder()
			h.GetMe(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_GetPurchases(t *testing.T) {
	type testCase struct {
		name        string
		ctxValue    int
		withContext bool
		mockSetup   func(*mock.MockRegistry)
		wantStatus  int
	}

	tests := []testCase{
		{
			name:        "Success",
			ctxValue:    1,
			withContext: true,
			mockSetup: func(r *mock.MockRegistry) {
				r.GetBuyerPurchasesUC = &mock.MockGetBuyerPurchasesUseCase{
					ExecuteFunc: func(_ context.Context, buyerID int) ([]model.Purchase, error) {
						if buyerID != 1 {
							return nil, errors.New("wrong ID")
						}
						return []model.Purchase{{ID: 1, ItemID: 1, BuyerID: 1}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Unauthorized_NoContext",
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := buyer.NewBuyerHandler(mockReg)

			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/purchases", nil)
			if tc.withContext {
				req = withBuyerID(req, tc.ctxValue)
			}

			w := httptest.NewRecorder()

			h.GetPurchases(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_GetAuctions(t *testing.T) {
	type testCase struct {
		name        string
		ctxValue    int
		withContext bool
		mockSetup   func(*mock.MockRegistry)
		wantStatus  int
	}
	tests := []testCase{
		{
			name:        "Success",
			ctxValue:    1,
			withContext: true,
			mockSetup: func(r *mock.MockRegistry) {
				r.GetBuyerAuctionsUC = &mock.MockGetBuyerAuctionsUseCase{
					ExecuteFunc: func(_ context.Context, _ int) ([]model.Auction, error) {
						now := time.Now()
						return []model.Auction{{ID: 1, Period: model.NewAuctionPeriod(now, &now, &now)}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Unauthorized_NoContext",
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := buyer.NewBuyerHandler(mockReg)
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/auctions", nil)
			if tc.withContext {
				req = withBuyerID(req, tc.ctxValue)
			}

			w := httptest.NewRecorder()
			h.GetAuctions(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_UpdatePassword(t *testing.T) {
	type testCase struct {
		name        string
		withContext bool
		body        any
		mockSetup   func(*mock.MockRegistry)
		wantStatus  int
	}

	tests := []testCase{
		{
			name:        "Success",
			withContext: true,
			body:        request.UpdatePassword{CurrentPassword: "old", NewPassword: "new"},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateBuyerPasswordUC = &mock.MockBuyerUpdatePasswordUseCase{
					ExecuteFunc: func(_ context.Context, _ int, _, _ string) error { return nil },
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "InvalidJSON",
			withContext: true,
			body:        "invalid",
			mockSetup:   func(_ *mock.MockRegistry) {},
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name:       "Unauthorized_NoContext",
			body:       request.UpdatePassword{},
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := buyer.NewBuyerHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}
			req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/password", bytes.NewReader(reqBody))
			if tc.withContext {
				req = withBuyerID(req, 1)
			}

			w := httptest.NewRecorder()
			h.UpdatePassword(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}
