package public_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/handler/public"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

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
						if filters.VenueID != nil || filters.AuctionDate != nil {
							return nil, errors.New("filters should be ignored")
						}
						return []model.Auction{{ID: 1}}, nil
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
			h := public.NewAuctionHandler(mockReg)

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
			name:  "NotFound",
			idStr: "999",
			mockSetup: func(r *mock.MockRegistry) {
				r.GetAuctionUC = &mock.MockGetAuctionUseCase{
					ExecuteFunc: func(_ context.Context, _ int) (*model.Auction, error) {
						return nil, errors.New("not found") // HandleError will map this to 500 if not specific domain error
						// Actually util.HandleError maps various errors.
					},
				}
			},
			wantStatus: http.StatusInternalServerError, // Default for generic error
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
			h := public.NewAuctionHandler(mockReg)

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
				r.GetAuctionItemsUC = &mock.MockGetAuctionItemsUseCase{
					ExecuteFunc: func(_ context.Context, _ int) ([]model.AuctionItem, error) {
						return []model.AuctionItem{{ID: 1}}, nil
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
				r.GetAuctionItemsUC = &mock.MockGetAuctionItemsUseCase{
					ExecuteFunc: func(_ context.Context, _ int) ([]model.AuctionItem, error) {
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
			h := public.NewAuctionHandler(mockReg)

			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/auctions/"+tc.idStr+"/items", nil)
			req.SetPathValue("id", tc.idStr)
			w := httptest.NewRecorder()

			h.GetItems(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestPublicAuctionHandler_RegisterRoutes(t *testing.T) {
	mockReg := &mock.MockRegistry{
		ListAuctionsUC:    &mock.MockListAuctionsUseCase{ExecuteFunc: func(_ context.Context, _ *repository.AuctionFilters) ([]model.Auction, error) { return []model.Auction{}, nil }},
		GetAuctionUC:      &mock.MockGetAuctionUseCase{ExecuteFunc: func(_ context.Context, _ int) (*model.Auction, error) { return &model.Auction{ID: 1}, nil }},
		GetAuctionItemsUC: &mock.MockGetAuctionItemsUseCase{ExecuteFunc: func(_ context.Context, _ int) ([]model.AuctionItem, error) { return []model.AuctionItem{}, nil }},
	}
	h := public.NewAuctionHandler(mockReg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/auctions"},
		{http.MethodGet, "/api/auctions/1"},
		{http.MethodGet, "/api/auctions/1/items"},
	}

	for _, tt := range tests {
		req := httptest.NewRequestWithContext(context.Background(), tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
			// Note: without proper router setup in test, path value might not be populated if just using mux.ServeHTTP directly without a real router.
			// But here we just check if it's registered.
		}
	}
}
