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
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestBuyerHandler_Create(t *testing.T) {
	type testCase struct {
		name       string
		body       interface{}
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			body: dto.CreateBuyerRequest{
				Name:         "Buyer 1",
				Email:        "buyer@example.com",
				Password:     "password",
				Organization: "Org 1",
				ContactInfo:  "Contact 1",
			},
			mockSetup: func(r *mock.MockRegistry) {
				r.CreateBuyerUC = &mock.MockCreateBuyerUseCase{
					ExecuteFunc: func(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error) {
						return &model.Buyer{ID: 1, Name: name, Organization: organization}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidJSON",
			body:       "invalid-json",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "UseCaseError",
			body: dto.CreateBuyerRequest{Email: "buyer@example.com"},
			mockSetup: func(r *mock.MockRegistry) {
				r.CreateBuyerUC = &mock.MockCreateBuyerUseCase{
					ExecuteFunc: func(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error) {
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
			h := handler.NewBuyerHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/buyers", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_Login(t *testing.T) {
	type testCase struct {
		name         string
		body         interface{}
		mockSetup    func(*mock.MockRegistry)
		wantStatus   int
		expectCookie bool
	}

	tests := []testCase{
		{
			name: "Success",
			body: dto.LoginBuyerRequest{Email: "buyer@example.com", Password: "password"},
			mockSetup: func(r *mock.MockRegistry) {
				r.LoginBuyerUC = &mock.MockLoginBuyerUseCase{
					ExecuteFunc: func(ctx context.Context, email, password string) (*model.Buyer, error) {
						return &model.Buyer{ID: 1, Name: "Buyer 1"}, nil
					},
				}
			},
			wantStatus:   http.StatusOK,
			expectCookie: true,
		},
		{
			name:       "InvalidJSON",
			body:       "invalid-json",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "InvalidCredentials",
			body: dto.LoginBuyerRequest{Email: "buyer@example.com", Password: "wrong"},
			mockSetup: func(r *mock.MockRegistry) {
				r.LoginBuyerUC = &mock.MockLoginBuyerUseCase{
					ExecuteFunc: func(ctx context.Context, email, password string) (*model.Buyer, error) {
						return nil, errors.New("invalid credentials")
					},
				}
			},
			wantStatus: http.StatusInternalServerError, // Handler just returns error, util maps to 500 usually
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := handler.NewBuyerHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/buyers/login", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			h.Login(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}

			if tc.expectCookie {
				cookies := w.Result().Cookies()
				found := false
				for _, c := range cookies {
					if c.Name == "buyer_session" && c.Value == "authenticated" {
						found = true
						break
					}
				}
				if !found {
					t.Error("expected buyer_session cookie")
				}
			}
		})
	}
}

func TestBuyerHandler_GetMyPurchases(t *testing.T) {
	type testCase struct {
		name       string
		cookies    map[string]string
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name:    "Success",
			cookies: map[string]string{"buyer_session": "authenticated", "buyer_id": "1"},
			mockSetup: func(r *mock.MockRegistry) {
				r.GetBuyerPurchasesUC = &mock.MockGetBuyerPurchasesUseCase{
					ExecuteFunc: func(ctx context.Context, buyerID int) ([]model.Purchase, error) {
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
			name:       "Unauthorized_NoSession",
			cookies:    map[string]string{},
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Unauthorized_NoID",
			cookies:    map[string]string{"buyer_session": "authenticated"},
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "InvalidIDFormat",
			cookies:    map[string]string{"buyer_session": "authenticated", "buyer_id": "invalid"},
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "UseCaseError",
			cookies: map[string]string{"buyer_session": "authenticated", "buyer_id": "1"},
			mockSetup: func(r *mock.MockRegistry) {
				r.GetBuyerPurchasesUC = &mock.MockGetBuyerPurchasesUseCase{
					ExecuteFunc: func(ctx context.Context, buyerID int) ([]model.Purchase, error) {
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
			h := handler.NewBuyerHandler(mockReg)

			req := httptest.NewRequest(http.MethodGet, "/api/buyers/me/purchases", nil)
			for k, v := range tc.cookies {
				req.AddCookie(&http.Cookie{Name: k, Value: v})
			}
			w := httptest.NewRecorder()

			h.GetMyPurchases(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_List(t *testing.T) {
	// Simple enough without table but following pattern
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
					ExecuteFunc: func(ctx context.Context) ([]model.Buyer, error) {
						return []model.Buyer{{Name: "B1"}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "UseCaseError",
			mockSetup: func(r *mock.MockRegistry) {
				r.ListBuyersUC = &mock.MockListBuyersUseCase{
					ExecuteFunc: func(ctx context.Context) ([]model.Buyer, error) {
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
			h := handler.NewBuyerHandler(mockReg)
			req := httptest.NewRequest(http.MethodGet, "/api/buyers", nil)
			w := httptest.NewRecorder()
			h.List(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_GetCurrentBuyer(t *testing.T) {
	// Reusing logic from GetMyPurchases table approach conceptually
	type testCase struct {
		name       string
		cookies    map[string]string
		wantStatus int
	}
	tests := []testCase{
		{name: "Success", cookies: map[string]string{"buyer_session": "authenticated", "buyer_id": "1"}, wantStatus: http.StatusOK},
		{name: "NoSession", cookies: map[string]string{}, wantStatus: http.StatusUnauthorized},
		{name: "NoID", cookies: map[string]string{"buyer_session": "authenticated"}, wantStatus: http.StatusUnauthorized},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := handler.NewBuyerHandler(&mock.MockRegistry{})
			req := httptest.NewRequest(http.MethodGet, "/api/buyers/me", nil)
			for k, v := range tc.cookies {
				req.AddCookie(&http.Cookie{Name: k, Value: v})
			}
			w := httptest.NewRecorder()
			h.GetCurrentBuyer(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_GetMyAuctions(t *testing.T) {
	type testCase struct {
		name       string
		cookies    map[string]string
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}
	// Logic is very similar to GetMyPurchases. Testing InvalidIDFormat specific to GetMyAuctions
	tests := []testCase{
		{
			name:    "Success",
			cookies: map[string]string{"buyer_session": "authenticated", "buyer_id": "1"},
			mockSetup: func(r *mock.MockRegistry) {
				r.GetBuyerAuctionsUC = &mock.MockGetBuyerAuctionsUseCase{
					ExecuteFunc: func(ctx context.Context, id int) ([]model.Auction, error) {
						now := time.Now()
						return []model.Auction{{ID: 1, StartTime: &now, EndTime: &now}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidIDFormat",
			cookies:    map[string]string{"buyer_session": "authenticated", "buyer_id": "invalid"},
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := handler.NewBuyerHandler(mockReg)
			req := httptest.NewRequest(http.MethodGet, "/api/buyers/me/auctions", nil)
			for k, v := range tc.cookies {
				req.AddCookie(&http.Cookie{Name: k, Value: v})
			}
			w := httptest.NewRecorder()
			h.GetMyAuctions(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_UpdatePassword(t *testing.T) {
	type testCase struct {
		name       string
		cookies    map[string]string
		body       interface{}
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name:    "Success",
			cookies: map[string]string{"buyer_session": "authenticated", "buyer_id": "1"},
			body:    dto.UpdatePasswordRequest{CurrentPassword: "old", NewPassword: "new"},
			mockSetup: func(r *mock.MockRegistry) {
				r.UpdateBuyerPasswordUC = &mock.MockBuyerUpdatePasswordUseCase{
					ExecuteFunc: func(ctx context.Context, id int, c, n string) error { return nil },
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidJSON",
			cookies:    map[string]string{"buyer_session": "authenticated", "buyer_id": "1"},
			body:       "invalid",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "InvalidIDFormat",
			cookies:    map[string]string{"buyer_session": "authenticated", "buyer_id": "invalid"},
			body:       dto.UpdatePasswordRequest{},
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := handler.NewBuyerHandler(mockReg)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}
			req := httptest.NewRequest(http.MethodPut, "/api/buyers/me/password", bytes.NewReader(reqBody))
			for k, v := range tc.cookies {
				req.AddCookie(&http.Cookie{Name: k, Value: v})
			}
			w := httptest.NewRecorder()
			h.UpdatePassword(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_Logout(t *testing.T) {
	// Simple
	h := handler.NewBuyerHandler(&mock.MockRegistry{})
	req := httptest.NewRequest(http.MethodPost, "/api/buyers/logout", nil)
	w := httptest.NewRecorder()
	h.Logout(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestBuyerHandler_RegisterRoutes(t *testing.T) {
	type testCase struct {
		name       string
		method     string
		path       string
		wantStatus int
	}

	tests := []testCase{
		{name: "Create_Post", method: http.MethodPost, path: "/api/buyers", wantStatus: http.StatusMethodNotAllowed},
		{name: "List_Get", method: http.MethodGet, path: "/api/buyers", wantStatus: http.StatusOK},
		{name: "Create_Put", method: http.MethodPut, path: "/api/buyers", wantStatus: http.StatusMethodNotAllowed},

		{name: "Login_Post", method: http.MethodPost, path: "/api/buyers/login", wantStatus: http.StatusOK},
		{name: "Login_Get", method: http.MethodGet, path: "/api/buyers/login", wantStatus: http.StatusMethodNotAllowed},

		{name: "Logout_Post", method: http.MethodPost, path: "/api/buyers/logout", wantStatus: http.StatusOK},
		{name: "Logout_Get", method: http.MethodGet, path: "/api/buyers/logout", wantStatus: http.StatusMethodNotAllowed},

		{name: "Me_Get", method: http.MethodGet, path: "/api/buyers/me", wantStatus: http.StatusOK},
		{name: "Me_Post", method: http.MethodPost, path: "/api/buyers/me", wantStatus: http.StatusMethodNotAllowed},

		{name: "Purchases_Get", method: http.MethodGet, path: "/api/buyers/me/purchases", wantStatus: http.StatusOK},
		{name: "Purchases_Post", method: http.MethodPost, path: "/api/buyers/me/purchases", wantStatus: http.StatusMethodNotAllowed},

		{name: "Auctions_Get", method: http.MethodGet, path: "/api/buyers/me/auctions", wantStatus: http.StatusOK},

		{name: "Password_Put", method: http.MethodPut, path: "/api/buyers/password", wantStatus: http.StatusOK},
		{name: "Password_Post", method: http.MethodPost, path: "/api/buyers/password", wantStatus: http.StatusMethodNotAllowed},
	}

	mockReg := &mock.MockRegistry{
		// Mock all UCs
		CreateBuyerUC:         &mock.MockCreateBuyerUseCase{ExecuteFunc: func(ctx context.Context, n, e, p, o, c string) (*model.Buyer, error) { return &model.Buyer{ID: 1}, nil }},
		ListBuyersUC:          &mock.MockListBuyersUseCase{ExecuteFunc: func(ctx context.Context) ([]model.Buyer, error) { return []model.Buyer{}, nil }},
		LoginBuyerUC:          &mock.MockLoginBuyerUseCase{ExecuteFunc: func(ctx context.Context, e, p string) (*model.Buyer, error) { return &model.Buyer{ID: 1}, nil }},
		GetBuyerPurchasesUC:   &mock.MockGetBuyerPurchasesUseCase{ExecuteFunc: func(ctx context.Context, id int) ([]model.Purchase, error) { return []model.Purchase{}, nil }},
		GetBuyerAuctionsUC:    &mock.MockGetBuyerAuctionsUseCase{ExecuteFunc: func(ctx context.Context, id int) ([]model.Auction, error) { return []model.Auction{}, nil }},
		UpdateBuyerPasswordUC: &mock.MockBuyerUpdatePasswordUseCase{ExecuteFunc: func(ctx context.Context, id int, c, n string) error { return nil }},
	}

	h := handler.NewBuyerHandler(mockReg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(body))
			// Add cookies for handlers that need them to pass auth check, to allow reaching method check or status 200
			req.AddCookie(&http.Cookie{Name: "buyer_session", Value: "authenticated"})
			req.AddCookie(&http.Cookie{Name: "buyer_id", Value: "1"})

			// Handle Bodies for strict JSON decoders if method matches
			if tc.method == http.MethodPost || tc.method == http.MethodPut {
				if tc.path == "/api/buyers" {
					body, _ = json.Marshal(dto.CreateBuyerRequest{})
				}
				if tc.path == "/api/buyers/login" {
					body, _ = json.Marshal(dto.LoginBuyerRequest{})
				}
				if tc.path == "/api/buyers/password" {
					body, _ = json.Marshal(dto.UpdatePasswordRequest{})
				}
				req = httptest.NewRequest(tc.method, tc.path, bytes.NewReader(body))
				// Re-add cookies
				req.AddCookie(&http.Cookie{Name: "buyer_session", Value: "authenticated"})
				req.AddCookie(&http.Cookie{Name: "buyer_id", Value: "1"})
			}

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}
