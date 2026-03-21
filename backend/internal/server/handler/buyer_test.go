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
	"github.com/seka/fish-auction/backend/internal/server/middleware"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func withBuyerID(req *http.Request, buyerID interface{}) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), middleware.BuyerIDKey, buyerID))
}

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
			h := handler.NewBuyerHandler(mockReg, &mock.MockSessionRepository{})

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/admin/buyers", bytes.NewReader(reqBody))
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
			sessionRepo := &mock.MockSessionRepository{NextSessionID: "buyer-session-1"}
			h := handler.NewBuyerHandler(mockReg, sessionRepo)

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
					if c.Name == "buyer_session" && c.Value == "buyer-session-1" {
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
		ctxValue    interface{}
		withContext bool
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name:    "Success",
			ctxValue: 1,
			withContext: true,
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
			name:       "Unauthorized_NoContext",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:    "UseCaseError",
			ctxValue: 1,
			withContext: true,
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
			h := handler.NewBuyerHandler(mockReg, &mock.MockSessionRepository{})

			req := httptest.NewRequest(http.MethodGet, "/api/buyers/me/purchases", nil)
			if tc.withContext {
				req = withBuyerID(req, tc.ctxValue)
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
			h := handler.NewBuyerHandler(mockReg, &mock.MockSessionRepository{})
			req := httptest.NewRequest(http.MethodGet, "/api/admin/buyers", nil)
			w := httptest.NewRecorder()
			h.List(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("expected %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_GetCurrentBuyer(t *testing.T) {
	type testCase struct {
		name       string
		withContext bool
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}
	tests := []testCase{
		{
			name:    "Success",
			withContext: true,
			mockSetup: func(r *mock.MockRegistry) {
				r.GetBuyerUC = &mock.MockGetBuyerUseCase{
					ExecuteFunc: func(ctx context.Context, id int) (*model.Buyer, error) {
						return &model.Buyer{ID: 1, Name: "Buyer 1"}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "NoContext",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			if tc.mockSetup != nil {
				tc.mockSetup(mockReg)
			}
			h := handler.NewBuyerHandler(mockReg, &mock.MockSessionRepository{})
			req := httptest.NewRequest(http.MethodGet, "/api/buyers/me", nil)
			if tc.withContext {
				req = withBuyerID(req, 1)
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
		ctxValue    interface{}
		withContext bool
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}
	tests := []testCase{
		{
			name:    "Success",
			ctxValue: 1,
			withContext: true,
			mockSetup: func(r *mock.MockRegistry) {
				r.GetBuyerAuctionsUC = &mock.MockGetBuyerAuctionsUseCase{
					ExecuteFunc: func(ctx context.Context, id int) ([]model.Auction, error) {
						now := time.Now()
						return []model.Auction{{ID: 1, Period: model.NewAuctionPeriod(now, &now, &now)}}, nil
					},
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Unauthorized_NoContext",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := handler.NewBuyerHandler(mockReg, &mock.MockSessionRepository{})
			req := httptest.NewRequest(http.MethodGet, "/api/buyers/me/auctions", nil)
			if tc.withContext {
				req = withBuyerID(req, tc.ctxValue)
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
		withContext bool
		body       interface{}
		mockSetup  func(*mock.MockRegistry)
		wantStatus int
	}

	tests := []testCase{
		{
			name:    "Success",
			withContext: true,
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
			withContext: true,
			body:       "invalid",
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Unauthorized_NoContext",
			body:       dto.UpdatePasswordRequest{},
			mockSetup:  func(r *mock.MockRegistry) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockReg := &mock.MockRegistry{}
			tc.mockSetup(mockReg)
			h := handler.NewBuyerHandler(mockReg, &mock.MockSessionRepository{})

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}
			req := httptest.NewRequest(http.MethodPut, "/api/buyers/me/password", bytes.NewReader(reqBody))
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

func TestBuyerHandler_Logout(t *testing.T) {
	sessionRepo := &mock.MockSessionRepository{
		Sessions: map[string]*model.Session{
			"buyer-session-1": {ID: "buyer-session-1", UserID: 1, Role: model.SessionRoleBuyer},
		},
	}
	h := handler.NewBuyerHandler(&mock.MockRegistry{}, sessionRepo)
	req := httptest.NewRequest(http.MethodPost, "/api/buyers/logout", nil)
	req.AddCookie(&http.Cookie{Name: "buyer_session", Value: "buyer-session-1"})
	w := httptest.NewRecorder()
	h.Logout(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if len(sessionRepo.DeletedSessionIDs) != 1 || sessionRepo.DeletedSessionIDs[0] != "buyer-session-1" {
		t.Errorf("expected buyer-session-1 to be deleted, got %#v", sessionRepo.DeletedSessionIDs)
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
		GetBuyerUC:            &mock.MockGetBuyerUseCase{ExecuteFunc: func(ctx context.Context, id int) (*model.Buyer, error) { return &model.Buyer{Name: "B1"}, nil }},
	}

	sessionRepo := &mock.MockSessionRepository{NextSessionID: "buyer-session-1"}
	h := handler.NewBuyerHandler(mockReg, sessionRepo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(body))

			// Handle Bodies for strict JSON decoders if method matches
			if tc.method == http.MethodPost || tc.method == http.MethodPut {
				if tc.path == "/api/buyers/login" {
					body, _ = json.Marshal(dto.LoginBuyerRequest{})
				}
				if tc.path == "/api/buyers/password" {
					body, _ = json.Marshal(dto.UpdatePasswordRequest{})
				}
				req = httptest.NewRequest(tc.method, tc.path, bytes.NewReader(body))
			}
			if tc.path == "/api/buyers/logout" {
				req.AddCookie(&http.Cookie{Name: "buyer_session", Value: "buyer-session-1"})
			}
			if tc.path != "/api/buyers/login" && tc.path != "/api/buyers/logout" {
				req = withBuyerID(req, 1)
			}

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}
