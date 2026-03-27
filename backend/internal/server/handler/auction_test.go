package handler_test

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
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
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
			body: dto.CreateAuctionRequest{
				VenueID:     1,
				AuctionDate: "2023-01-01",
				Status:      "Scheduled",
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
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidJSON",
			body:       "invalid-json",
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "InvalidDateFormat",
			body: dto.CreateAuctionRequest{
				VenueID:     1,
				AuctionDate: "invalid-date",
			},
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCaseError",
			body: dto.CreateAuctionRequest{
				VenueID:     1,
				AuctionDate: "2023-01-01",
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
			h := handler.NewAdminAuctionHandler(mockReg)

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

func TestPublicAuctionHandler_List(t *testing.T) {
	type testCase struct {
		name        string
		queryParams map[string]string
		mockSetup   func(*mock.MockRegistry)
		wantStatus  int
	}

	tests := []testCase{
		{
			name: "Success",
			mockSetup: func(r *mock.MockRegistry) {
				r.ListAuctionsUC = &mock.MockListAuctionsUseCase{
					ExecuteFunc: func(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) {
						return []model.Auction{{ID: 1}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "WithFilters",
			queryParams: map[string]string{
				"venue_id": "1",
				"date":     "2023-01-01",
				"status":   "Scheduled",
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.ListAuctionsUC = &mock.MockListAuctionsUseCase{
					ExecuteFunc: func(_ context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
						if filters.VenueID == nil || *filters.VenueID != 1 {
							return nil, errors.New("filter mismatch")
						}
						return []model.Auction{{ID: 1}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "InvalidFilters_Ignored",
			queryParams: map[string]string{
				"venue_id": "invalid",
				"date":     "invalid",
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.ListAuctionsUC = &mock.MockListAuctionsUseCase{
					ExecuteFunc: func(_ context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
						if filters.VenueID != nil {
							return nil, errors.New("expected nil VenueID")
						}
						return []model.Auction{}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "UseCaseError",
			mockSetup: func(r *mock.MockRegistry) {
				r.ListAuctionsUC = &mock.MockListAuctionsUseCase{
					ExecuteFunc: func(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) {
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
			h := handler.NewPublicAuctionHandler(mockReg)

			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/auctions", nil)
			q := req.URL.Query()
			for k, v := range tc.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			w := httptest.NewRecorder()

			h.List(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestPublicAuctionHandler_Get(t *testing.T) {
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
				r.GetAuctionUC = &mock.MockGetAuctionUseCase{
					ExecuteFunc: func(_ context.Context, _ int) (*model.Auction, error) {
						return &model.Auction{ID: 1}, nil
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
				r.GetAuctionUC = &mock.MockGetAuctionUseCase{
					ExecuteFunc: func(_ context.Context, _ int) (*model.Auction, error) {
						return nil, errors.New("not found")
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
			h := handler.NewPublicAuctionHandler(mockReg)

			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/auctions/"+tc.idStr, nil)
			req.SetPathValue("id", tc.idStr)
			w := httptest.NewRecorder()

			h.Get(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestPublicAuctionHandler_GetItems(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockReg := &mock.MockRegistry{
			GetAuctionItemsUC: &mock.MockGetAuctionItemsUseCase{
				ExecuteFunc: func(_ context.Context, _ int) ([]model.AuctionItem, error) {
					return []model.AuctionItem{}, nil
				},
			},
		}
		h := handler.NewPublicAuctionHandler(mockReg)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/auctions/1/items", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.GetItems(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		h := handler.NewPublicAuctionHandler(&mock.MockRegistry{})
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/auctions/invalid/items", nil)
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()
		h.GetItems(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
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
			body: dto.UpdateAuctionRequest{
				VenueID:     1,
				AuctionDate: "2023-01-01",
				Status:      "Completed",
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
			name:       "InvalidID",
			idStr:      "invalid",
			body:       dto.UpdateAuctionRequest{},
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "InvalidJSON",
			idStr:      "1",
			body:       "invalid-json",
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:  "UseCaseError",
			idStr: "1",
			body: dto.UpdateAuctionRequest{
				AuctionDate: "2023-01-01",
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
			h := handler.NewAdminAuctionHandler(mockReg)

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
			body:  dto.UpdateAuctionStatusRequest{Status: "Closed"},
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
			name:       "InvalidID",
			idStr:      "invalid",
			body:       dto.UpdateAuctionStatusRequest{Status: "Closed"},
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "UseCaseError",
			idStr: "1",
			body:  dto.UpdateAuctionStatusRequest{Status: "Closed"},
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
			h := handler.NewAdminAuctionHandler(mockReg)

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
			h := handler.NewAdminAuctionHandler(mockReg)

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
	}

	t.Run("Public", func(t *testing.T) {
		h := handler.NewPublicAuctionHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/auctions", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Admin", func(t *testing.T) {
		h := handler.NewAdminAuctionHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auctions", strings.NewReader(`{"auction_date":"2023-01-01"}`))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}
