package handler_test

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
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAuctionHandler_Create(t *testing.T) {
	type testCase struct {
		name       string
		body       interface{} // string or struct
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
					ExecuteFunc: func(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
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
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError, // util.HandleError returns 500 for generic errors or maybe bad request? Check util implementation.
			// Checking util.HandleError implementation usually maps json decode error to 500 unless handled specifically.
			// Let's assume 500 for now based on previous simple HandleError usage.
			// Actually standard json.Decode error is just an error.
		},
		{
			name: "InvalidDateFormat",
			body: dto.CreateAuctionRequest{
				VenueID:     1,
				AuctionDate: "invalid-date",
			},
			mockSetup:  func(r *mock.MockRegistry) {},
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
					ExecuteFunc: func(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
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
			h := handler.NewAuctionHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/auctions", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAuctionHandler_List(t *testing.T) {
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
					ExecuteFunc: func(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
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
					ExecuteFunc: func(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
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
					ExecuteFunc: func(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
						// Handler ignores invalid parse, so filters should be nil/empty
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
					ExecuteFunc: func(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
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
			h := handler.NewAuctionHandler(mockReg)

			req := httptest.NewRequest(http.MethodGet, "/api/auctions", nil)
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

func TestAuctionHandler_Get(t *testing.T) {
	type testCase struct {
		name       string
		path       string
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			path: "/api/auctions/1",
			mockSetup: func(r *mock.MockRegistry) {
				r.GetAuctionUC = &mock.MockGetAuctionUseCase{
					ExecuteFunc: func(ctx context.Context, id int) (*model.Auction, error) {
						return &model.Auction{ID: 1}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "InvalidID",
			path: "/api/auctions/invalid",
			mockSetup: func(r *mock.MockRegistry) {
				// No UC call expected
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "NotFound",
			path: "/api/auctions/999",
			mockSetup: func(r *mock.MockRegistry) {
				r.GetAuctionUC = &mock.MockGetAuctionUseCase{
					ExecuteFunc: func(ctx context.Context, id int) (*model.Auction, error) {
						return nil, errors.New("not found")
					},
				}
			},
			wantStatus: http.StatusInternalServerError, // Handler maps error to 500 via util.HandleError
		},
		// Special suffixes routed to other methods
		{
			name: "RouteToItems",
			path: "/api/auctions/1/items",
			mockSetup: func(r *mock.MockRegistry) {
				r.GetAuctionItemsUC = &mock.MockGetAuctionItemsUseCase{
					ExecuteFunc: func(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
						return []model.AuctionItem{}, nil
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
			h := handler.NewAuctionHandler(mockReg)

			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			w := httptest.NewRecorder()

			h.Get(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAuctionHandler_GetItems(t *testing.T) {
	// Simple test without table as logic is similar to Get
	// But let's add InvalidID case
	t.Run("Success", func(t *testing.T) {
		mockReg := &mock.MockRegistry{
			GetAuctionItemsUC: &mock.MockGetAuctionItemsUseCase{
				ExecuteFunc: func(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
					return []model.AuctionItem{}, nil
				},
			},
		}
		h := handler.NewAuctionHandler(mockReg)
		req := httptest.NewRequest(http.MethodGet, "/api/auctions/1/items", nil)
		w := httptest.NewRecorder()
		h.GetItems(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		h := handler.NewAuctionHandler(&mock.MockRegistry{})
		req := httptest.NewRequest(http.MethodGet, "/api/auctions/invalid/items", nil)
		w := httptest.NewRecorder()
		h.GetItems(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockReg := &mock.MockRegistry{
			GetAuctionItemsUC: &mock.MockGetAuctionItemsUseCase{
				ExecuteFunc: func(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
					return nil, errors.New("db error")
				},
			},
		}
		h := handler.NewAuctionHandler(mockReg)
		req := httptest.NewRequest(http.MethodGet, "/api/auctions/1/items", nil)
		w := httptest.NewRecorder()
		h.GetItems(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestAuctionHandler_Update(t *testing.T) {
	type testCase struct {
		name       string
		path       string
		body       interface{}
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			path: "/api/auctions/1",
			body: dto.UpdateAuctionRequest{
				VenueID:     1,
				AuctionDate: "2023-01-01",
				Status:      "Completed",
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateAuctionUC = &mock.MockUpdateAuctionUseCase{
					ExecuteFunc: func(ctx context.Context, auction *model.Auction) error {
						return nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidID",
			path:       "/api/auctions/invalid",
			body:       dto.UpdateAuctionRequest{},
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "InvalidJSON",
			path:       "/api/auctions/1",
			body:       "invalid-json",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "InvalidDate",
			path: "/api/auctions/1",
			body: dto.UpdateAuctionRequest{
				AuctionDate: "invalid",
			},
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCaseError",
			path: "/api/auctions/1",
			body: dto.UpdateAuctionRequest{
				AuctionDate: "2023-01-01",
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateAuctionUC = &mock.MockUpdateAuctionUseCase{
					ExecuteFunc: func(ctx context.Context, auction *model.Auction) error {
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
			h := handler.NewAuctionHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequest(http.MethodPut, tc.path, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			h.Update(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAuctionHandler_UpdateStatus(t *testing.T) {
	type testCase struct {
		name       string
		path       string
		body       interface{}
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			path: "/api/auctions/1/status",
			body: dto.UpdateAuctionStatusRequest{Status: "Closed"},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateAuctionStatusUC = &mock.MockUpdateAuctionStatusUseCase{
					ExecuteFunc: func(ctx context.Context, id int, status model.AuctionStatus) error {
						return nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidID",
			path:       "/api/auctions/invalid/status",
			body:       dto.UpdateAuctionStatusRequest{Status: "Closed"},
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "InvalidJSON",
			path:       "/api/auctions/1/status",
			body:       "invalid-json",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "UseCaseError",
			path: "/api/auctions/1/status",
			body: dto.UpdateAuctionStatusRequest{Status: "Closed"},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateAuctionStatusUC = &mock.MockUpdateAuctionStatusUseCase{
					ExecuteFunc: func(ctx context.Context, id int, status model.AuctionStatus) error {
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
			h := handler.NewAuctionHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequest(http.MethodPatch, tc.path, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			h.UpdateStatus(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAuctionHandler_Delete(t *testing.T) {
	type testCase struct {
		name       string
		path       string
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			path: "/api/auctions/1",
			mockSetup: func(r *mock.MockRegistry) {
				r.DeleteAuctionUC = &mock.MockDeleteAuctionUseCase{
					ExecuteFunc: func(ctx context.Context, id int) error {
						return nil
					},
				}
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "InvalidID",
			path:       "/api/auctions/invalid",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCaseError",
			path: "/api/auctions/1",
			mockSetup: func(r *mock.MockRegistry) {
				r.DeleteAuctionUC = &mock.MockDeleteAuctionUseCase{
					ExecuteFunc: func(ctx context.Context, id int) error {
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
			h := handler.NewAuctionHandler(mockReg)

			req := httptest.NewRequest(http.MethodDelete, tc.path, nil)
			w := httptest.NewRecorder()

			h.Delete(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestAuctionHandler_RegisterRoutes(t *testing.T) {
	type testCase struct {
		name       string
		method     string
		path       string
		wantStatus int // 405 MethodNotAllowed or 200 OK (if mocked to success)
	}

	tests := []testCase{
		{name: "Create_Post", method: http.MethodPost, path: "/api/auctions", wantStatus: http.StatusOK},
		{name: "Create_Put", method: http.MethodPut, path: "/api/auctions", wantStatus: http.StatusMethodNotAllowed},

		{name: "List_Get", method: http.MethodGet, path: "/api/auctions", wantStatus: http.StatusOK},
		// Note: /api/auctions matches both POST and GET.

		{name: "Get_Get", method: http.MethodGet, path: "/api/auctions/1", wantStatus: http.StatusOK},
		{name: "Update_Put", method: http.MethodPut, path: "/api/auctions/1", wantStatus: http.StatusOK},
		{name: "Delete_Delete", method: http.MethodDelete, path: "/api/auctions/1", wantStatus: http.StatusNoContent},

		{name: "GetItems_Get", method: http.MethodGet, path: "/api/auctions/1/items", wantStatus: http.StatusOK},
		{name: "GetItems_Post", method: http.MethodPost, path: "/api/auctions/1/items", wantStatus: http.StatusMethodNotAllowed}, // Fallthrough to switch default in suffix handler? Or handled? Handler detail: it checks GetItems first.

		{name: "UpdateStatus_Patch", method: http.MethodPatch, path: "/api/auctions/1/status", wantStatus: http.StatusOK},
		{name: "UpdateStatus_Post", method: http.MethodPost, path: "/api/auctions/1/status", wantStatus: http.StatusMethodNotAllowed},

		{name: "InvalidMethod_Post_Detail", method: http.MethodPost, path: "/api/auctions/1", wantStatus: http.StatusMethodNotAllowed},
	}

	mockReg := &mock.MockRegistry{
		// Mock all UCs to return success for routing tests
		CreateAuctionUC: &mock.MockCreateAuctionUseCase{ExecuteFunc: func(ctx context.Context, a *model.Auction) (*model.Auction, error) { a.ID = 1; return a, nil }},
		ListAuctionsUC: &mock.MockListAuctionsUseCase{ExecuteFunc: func(ctx context.Context, f *repository.AuctionFilters) ([]model.Auction, error) {
			return []model.Auction{}, nil
		}},
		GetAuctionUC:          &mock.MockGetAuctionUseCase{ExecuteFunc: func(ctx context.Context, id int) (*model.Auction, error) { return &model.Auction{ID: 1}, nil }},
		UpdateAuctionUC:       &mock.MockUpdateAuctionUseCase{ExecuteFunc: func(ctx context.Context, a *model.Auction) error { return nil }},
		DeleteAuctionUC:       &mock.MockDeleteAuctionUseCase{ExecuteFunc: func(ctx context.Context, id int) error { return nil }},
		GetAuctionItemsUC:     &mock.MockGetAuctionItemsUseCase{ExecuteFunc: func(ctx context.Context, id int) ([]model.AuctionItem, error) { return []model.AuctionItem{}, nil }},
		UpdateAuctionStatusUC: &mock.MockUpdateAuctionStatusUseCase{ExecuteFunc: func(ctx context.Context, id int, s model.AuctionStatus) error { return nil }},
	}

	h := handler.NewAuctionHandler(mockReg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			if tc.method == http.MethodPost || tc.method == http.MethodPut || tc.method == http.MethodPatch {
				// Provide minimal valid body to avoid 400/500 from Handler before method check if check is late?
				// Actually Handler checks method first usually in RegisterRoutes.
				// But to succeed we need body for 200 cases.
				if tc.path == "/api/auctions" && tc.method == http.MethodPost {
					body, _ = json.Marshal(dto.CreateAuctionRequest{AuctionDate: "2023-01-01"})
				} else if tc.path == "/api/auctions/1" && tc.method == http.MethodPut {
					body, _ = json.Marshal(dto.UpdateAuctionRequest{AuctionDate: "2023-01-01"})
				} else if tc.path == "/api/auctions/1/status" && tc.method == http.MethodPatch {
					body, _ = json.Marshal(dto.UpdateAuctionStatusRequest{Status: "Closed"})
				} else {
					body = []byte("{}")
				}
			}

			req := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(body))
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}
