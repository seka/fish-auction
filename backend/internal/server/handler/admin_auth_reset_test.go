package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAdminAuthResetHandler(t *testing.T) {
	t.Run("RequestReset_Success", func(t *testing.T) {
		mockReqUC := &mock.MockAdminRequestPasswordResetUseCase{
			ExecuteFunc: func(_ context.Context, _ string) error {
				return nil
			},
		}
		mockReg := &mock.MockRegistry{RequestAdminPasswordResetUC: mockReqUC}
		h := handler.NewAdminAuthResetHandler(mockReg)

		reqBody := map[string]string{"email": "admin@example.com"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/admin/password-reset/request", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.RequestReset(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("VerifyToken_Success", func(t *testing.T) {
		mockVerifyUC := &mock.MockAdminVerifyResetTokenUseCase{
			ExecuteFunc: func(_ context.Context, token string) error {
				if token == "valid-token" {
					return nil
				}
				return &errors.UnauthorizedError{Message: "Invalid or expired token"}
			},
		}
		mockReg := &mock.MockRegistry{VerifyAdminResetTokenUC: mockVerifyUC}
		h := handler.NewAdminAuthResetHandler(mockReg)

		reqBody := map[string]string{"token": "valid-token"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/admin/password-reset/verify", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.VerifyToken(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("VerifyToken_Unauthorized", func(t *testing.T) {
		mockVerifyUC := &mock.MockAdminVerifyResetTokenUseCase{
			ExecuteFunc: func(_ context.Context, _ string) error {
				return &errors.UnauthorizedError{Message: "Invalid or expired token"}
			},
		}
		mockReg := &mock.MockRegistry{VerifyAdminResetTokenUC: mockVerifyUC}
		h := handler.NewAdminAuthResetHandler(mockReg)

		reqBody := map[string]string{"token": "invalid-token"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/admin/password-reset/verify", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.VerifyToken(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("ConfirmReset_Success", func(t *testing.T) {
		mockResetUC := &mock.MockAdminResetPasswordUseCase{
			ExecuteFunc: func(_ context.Context, _, _ string) error {
				return nil
			},
		}
		mockReg := &mock.MockRegistry{ResetAdminPasswordUC: mockResetUC}
		h := handler.NewAdminAuthResetHandler(mockReg)

		reqBody := map[string]string{"token": "token123", "new_password": "newpass"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/admin/password-reset/confirm", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.ConfirmReset(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("RegisterRoutes_MethodNotAllowed", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := handler.NewAdminAuthResetHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/admin/password-reset/request", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}
