package admin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/request"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
	"github.com/seka/fish-auction/backend/internal/server/util"
)

func TestAdminAuctionHandler_Create(t *testing.T) {
	type testCase struct {
		name       string
		body       any
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			body: request.CreateAuction{
				VenueID: 1,
				Status:  "Scheduled",
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.CreateAuctionUC = &mock.MockCreateAuctionUseCase{
					ExecuteFunc: func(_ context.Context, auction *model.Auction) (*model.Auction, error) {
						auction.ID = 1
						auction.CreatedAt = time.Now()
						auction.UpdatedAt = time.Now()
						return auction, nil
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
			name: "InvalidStartAtFormat",
			body: map[string]any{
				"venue_id": 1,
				"start_at": "2026-03-15",
				"status":   "scheduled",
			},
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "InvalidEndAtFormat",
			body: map[string]any{
				"venue_id": 1,
				"start_at": "2026-03-15T09:00:00+09:00",
				"end_at":   "not-a-timestamp",
				"status":   "scheduled",
			},
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCaseError",
			body: request.CreateAuction{
				VenueID: 1,
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.CreateAuctionUC = &mock.MockCreateAuctionUseCase{
					ExecuteFunc: func(_ context.Context, _ *model.Auction) (*model.Auction, error) {
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
			h := admin.NewAuctionHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auctions", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAdminAuctionHandler_Update(t *testing.T) {
	type testCase struct {
		name       string
		idStr      string
		body       any
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name:  "Success",
			idStr: "1",
			body: request.UpdateAuction{
				VenueID: 1,
				Status:  "Completed",
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateAuctionUC = &mock.MockUpdateAuctionUseCase{
					ExecuteFunc: func(_ context.Context, _ *model.Auction) error {
						return nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidJSON",
			idStr:      "1",
			body:       "invalid-json",
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:  "InvalidStartAtFormat",
			idStr: "1",
			body: map[string]any{
				"venue_id": 1,
				"start_at": "2026-03-15",
				"status":   "scheduled",
			},
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "InvalidID",
			idStr:      "invalid",
			body:       request.UpdateAuction{},
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "UseCaseError",
			idStr: "1",
			body: request.UpdateAuction{
				VenueID: 1,
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateAuctionUC = &mock.MockUpdateAuctionUseCase{
					ExecuteFunc: func(_ context.Context, _ *model.Auction) error {
						return errors.New("db error")
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
			h := admin.NewAuctionHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/auctions/"+tc.idStr, bytes.NewReader(reqBody))
			req.SetPathValue("id", tc.idStr)
			w := httptest.NewRecorder()

			h.Update(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAdminAuctionHandler_UpdateStatus(t *testing.T) {
	type testCase struct {
		name       string
		idStr      string
		body       any
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name:  "Success",
			idStr: "1",
			body:  request.UpdateAuctionStatus{Status: "Closed"},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateAuctionStatusUC = &mock.MockUpdateAuctionStatusUseCase{
					ExecuteFunc: func(_ context.Context, _ int, _ model.AuctionStatus) error {
						return nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "InProgressSuccessWithStartAt",
			idStr: "1",
			body: map[string]any{
				"status":   "in_progress",
				"start_at": "2026-03-15T09:00:00+09:00",
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.GetAuctionUC = &mock.MockGetAuctionUseCase{
					ExecuteFunc: func(_ context.Context, _ int) (*model.Auction, error) {
						return &model.Auction{
							ID:      1,
							VenueID: 1,
							Status:  model.AuctionStatusScheduled,
							Period:  model.NewAuctionPeriod(nil, nil),
						}, nil
					},
				}
				r.UpdateAuctionUC = &mock.MockUpdateAuctionUseCase{
					ExecuteFunc: func(_ context.Context, _ *model.Auction) error {
						return nil
					},
				}
				r.UpdateAuctionStatusUC = &mock.MockUpdateAuctionStatusUseCase{
					ExecuteFunc: func(_ context.Context, _ int, _ model.AuctionStatus) error {
						return nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "InProgressWithoutStartAt",
			idStr: "1",
			body:  request.UpdateAuctionStatus{Status: "in_progress"},
			mockSetup: func(_ *mock.MockRegistry) {
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "InvalidStartAtFormat",
			idStr: "1",
			body: map[string]any{
				"status":   "in_progress",
				"start_at": "2026-03-15",
			},
			mockSetup: func(_ *mock.MockRegistry) {
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "InvalidID",
			idStr:      "invalid",
			body:       request.UpdateAuctionStatus{Status: "Closed"},
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "UseCaseError",
			idStr: "1",
			body:  request.UpdateAuctionStatus{Status: "Closed"},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateAuctionStatusUC = &mock.MockUpdateAuctionStatusUseCase{
					ExecuteFunc: func(_ context.Context, _ int, _ model.AuctionStatus) error {
						return errors.New("db error")
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
			h := admin.NewAuctionHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequestWithContext(context.Background(), http.MethodPatch, "/auctions/"+tc.idStr+"/status", bytes.NewReader(reqBody))
			req.SetPathValue("id", tc.idStr)
			w := httptest.NewRecorder()

			h.UpdateStatus(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAdminAuctionHandler_Delete(t *testing.T) {
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
				r.DeleteAuctionUC = &mock.MockDeleteAuctionUseCase{
					ExecuteFunc: func(_ context.Context, _ int) error {
						return nil
					},
				}
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "InvalidID",
			idStr:      "invalid",
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := admin.NewAuctionHandler(mockReg)

			req := httptest.NewRequestWithContext(context.Background(), http.MethodDelete, "/auctions/"+tc.idStr, nil)
			req.SetPathValue("id", tc.idStr)
			w := httptest.NewRecorder()

			h.Delete(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAuctionHandler_RegisterRoutes(t *testing.T) {
	mockReg := &mock.MockRegistry{
		CreateAuctionUC: &mock.MockCreateAuctionUseCase{ExecuteFunc: func(_ context.Context, a *model.Auction) (*model.Auction, error) { a.ID = 1; return a, nil }},
		ListAuctionsUC: &mock.MockListAuctionsUseCase{ExecuteFunc: func(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) {
			return []model.Auction{}, nil
		}},
		GetAuctionUC:          &mock.MockGetAuctionUseCase{ExecuteFunc: func(_ context.Context, _ int) (*model.Auction, error) { return &model.Auction{ID: 1}, nil }},
		UpdateAuctionUC:       &mock.MockUpdateAuctionUseCase{ExecuteFunc: func(_ context.Context, _ *model.Auction) error { return nil }},
		DeleteAuctionUC:       &mock.MockDeleteAuctionUseCase{ExecuteFunc: func(_ context.Context, _ int) error { return nil }},
		GetAuctionItemsUC:     &mock.MockGetAuctionItemsUseCase{ExecuteFunc: func(_ context.Context, _ int) ([]model.AuctionItem, error) { return []model.AuctionItem{}, nil }},
		UpdateAuctionStatusUC: &mock.MockUpdateAuctionStatusUseCase{ExecuteFunc: func(_ context.Context, _ int, _ model.AuctionStatus) error { return nil }},
		ReorderItemsUC:        &mock.MockReorderItemsUseCase{ExecuteFunc: func(_ context.Context, _ int, _ []int) error { return nil }},
	}

	h := admin.NewAuctionHandler(mockReg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Create
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auctions", strings.NewReader(`{}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Create: expected 201, got %d", w.Code)
	}

	// Reorder - Success (204)
	req = httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/auctions/1/reorder", strings.NewReader(`{"ids":[1,2,3]}`))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Reorder Success: expected 204, got %d", w.Code)
	}

	// Reorder - Invalid JSON (400)
	req = httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/auctions/1/reorder", strings.NewReader(`invalid json`))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Reorder Invalid JSON: expected 400, got %d", w.Code)
	}
	var errResp util.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Errorf("Reorder Invalid JSON: expected JSON error response, got error: %v", err)
	}
}
