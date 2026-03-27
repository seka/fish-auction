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

func withBuyerID(req *http.Request, buyerID any) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), middleware.BuyerIDKey, buyerID))
}

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
			body: dto.CreateBuyerRequest{
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
			wantStatus: http.StatusOK,
		},
		{
			name:       "InvalidJSON",
			body:       "invalid-json",
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "UseCaseError",
			body: dto.CreateBuyerRequest{Email: "buyer@example.com"},
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
			h := handler.NewAdminBuyerHandler(mockReg)

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

func TestBuyerHandler_Login(t *testing.T) {
	type testCase struct {
		name         string
		body         any
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
					ExecuteFunc: func(_ context.Context, _, _ string) (*model.Buyer, error) {
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
			mockSetup:  func(_ *mock.MockRegistry) {},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "InvalidCredentials",
			body: dto.LoginBuyerRequest{Email: "buyer@example.com", Password: "wrong"},
			mockSetup: func(r *mock.MockRegistry) {
				r.LoginBuyerUC = &mock.MockLoginBuyerUseCase{
					ExecuteFunc: func(_ context.Context, _, _ string) (*model.Buyer, error) {
						return nil, errors.New("invalid credentials")
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
			sessionRepo := &mock.MockSessionRepository{NextSessionID: "buyer-session-1"}
			h := handler.NewBuyerHandler(mockReg, sessionRepo)

			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/buyer/login", bytes.NewReader(reqBody))
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
		name        string
		ctxValue    any
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
		{
			name:        "UseCaseError",
			ctxValue:    1,
			withContext: true,
			mockSetup: func(r *mock.MockRegistry) {
				r.GetBuyerPurchasesUC = &mock.MockGetBuyerPurchasesUseCase{
					ExecuteFunc: func(_ context.Context, _ int) ([]model.Purchase, error) {
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

			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/me/purchases", nil)
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
		{
			name: "UseCaseError",
			mockSetup: func(r *mock.MockRegistry) {
				r.ListBuyersUC = &mock.MockListBuyersUseCase{
					ExecuteFunc: func(_ context.Context) ([]model.Buyer, error) {
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
			h := handler.NewAdminBuyerHandler(mockReg)
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/buyers", nil)
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
			name:       "NoContext",
			mockSetup:  func(_ *mock.MockRegistry) {},
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
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/me", nil)
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
		name        string
		ctxValue    any
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
			h := handler.NewBuyerHandler(mockReg, &mock.MockSessionRepository{})
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/me/auctions", nil)
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
			body:        dto.UpdatePasswordRequest{CurrentPassword: "old", NewPassword: "new"},
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
			body:       dto.UpdatePasswordRequest{},
			mockSetup:  func(_ *mock.MockRegistry) {},
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

func TestBuyerHandler_Logout(t *testing.T) {
	sessionRepo := &mock.MockSessionRepository{
		Sessions: map[string]*model.Session{
			"buyer-session-1": {ID: "buyer-session-1", UserID: 1, Role: model.SessionRoleBuyer},
		},
	}
	h := handler.NewBuyerHandler(&mock.MockRegistry{}, sessionRepo)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/buyer/logout", nil)
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

func TestBuyerHandler_RegisterPublicRoutes(t *testing.T) {
	type testCase struct {
		name       string
		method     string
		path       string
		wantStatus int
	}

	tests := []testCase{
		{name: "Login_Post", method: http.MethodPost, path: "/api/buyer/login", wantStatus: http.StatusOK},
		{name: "Login_Get", method: http.MethodGet, path: "/api/buyer/login", wantStatus: http.StatusMethodNotAllowed},
		{name: "Logout_Post", method: http.MethodPost, path: "/api/buyer/logout", wantStatus: http.StatusOK},
	}

	mockReg := &mock.MockRegistry{
		LoginBuyerUC: &mock.MockLoginBuyerUseCase{ExecuteFunc: func(_ context.Context, _, _ string) (*model.Buyer, error) { return &model.Buyer{ID: 1}, nil }},
	}

	sessionRepo := &mock.MockSessionRepository{NextSessionID: "buyer-session-1"}
	h := handler.NewBuyerHandler(mockReg, sessionRepo)
	mux := http.NewServeMux()
	h.RegisterPublicRoutes(mux)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			if tc.method == http.MethodPost {
				body, _ = json.Marshal(map[string]string{"email": "", "password": ""})
			}
			req := httptest.NewRequestWithContext(context.Background(), tc.method, tc.path, bytes.NewReader(body))
			if tc.path == "/api/buyer/logout" {
				req.AddCookie(&http.Cookie{Name: "buyer_session", Value: "buyer-session-1"})
			}

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}

func TestBuyerHandler_RegisterBuyerRoutes(t *testing.T) {
	type testCase struct {
		name        string
		method      string
		path        string
		body        []byte
		withContext bool
		wantStatus  int
	}

	tests := []testCase{
		{name: "Me_Get", method: http.MethodGet, path: "/me", withContext: true, wantStatus: http.StatusOK},
		{name: "Me_Post", method: http.MethodPost, path: "/me", wantStatus: http.StatusMethodNotAllowed},
		{name: "Purchases_Get", method: http.MethodGet, path: "/me/purchases", withContext: true, wantStatus: http.StatusOK},
		{name: "Purchases_Post", method: http.MethodPost, path: "/me/purchases", wantStatus: http.StatusMethodNotAllowed},
		{name: "Auctions_Get", method: http.MethodGet, path: "/me/auctions", withContext: true, wantStatus: http.StatusOK},
		{name: "Password_Put", method: http.MethodPut, path: "/password", body: []byte(`{"current_password":"old","new_password":"new"}`), withContext: true, wantStatus: http.StatusOK},
		{name: "Password_Post", method: http.MethodPost, path: "/password", wantStatus: http.StatusMethodNotAllowed},
	}

	mockReg := &mock.MockRegistry{
		GetBuyerUC: &mock.MockGetBuyerUseCase{
			ExecuteFunc: func(_ context.Context, _ int) (*model.Buyer, error) {
				return &model.Buyer{ID: 1, Name: "B1"}, nil
			},
		},
		GetBuyerPurchasesUC: &mock.MockGetBuyerPurchasesUseCase{
			ExecuteFunc: func(_ context.Context, _ int) ([]model.Purchase, error) {
				return []model.Purchase{}, nil
			},
		},
		GetBuyerAuctionsUC: &mock.MockGetBuyerAuctionsUseCase{
			ExecuteFunc: func(_ context.Context, _ int) ([]model.Auction, error) {
				return []model.Auction{}, nil
			},
		},
		UpdateBuyerPasswordUC: &mock.MockBuyerUpdatePasswordUseCase{
			ExecuteFunc: func(_ context.Context, _ int, _, _ string) error {
				return nil
			},
		},
	}

	h := handler.NewBuyerHandler(mockReg, &mock.MockSessionRepository{})
	mux := http.NewServeMux()
	h.RegisterBuyerRoutes(mux)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(context.Background(), tc.method, tc.path, bytes.NewReader(tc.body))
			if tc.withContext {
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
