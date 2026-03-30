package public_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/handler/public"
	"github.com/seka/fish-auction/backend/internal/server/handler/public/request"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestBuyerAuthHandler_Login(t *testing.T) {
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
			body: request.Login{Email: "buyer@example.com", Password: "password"},
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
			body: request.Login{Email: "buyer@example.com", Password: "wrong"},
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
			h := public.NewBuyerAuthHandler(mockReg, sessionRepo)

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

func TestBuyerAuthHandler_Logout(t *testing.T) {
	sessionRepo := &mock.MockSessionRepository{
		Sessions: map[string]*model.Session{
			"buyer-session-1": {ID: "buyer-session-1", UserID: 1, Role: model.SessionRoleBuyer},
		},
	}
	h := public.NewBuyerAuthHandler(&mock.MockRegistry{}, sessionRepo)
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
