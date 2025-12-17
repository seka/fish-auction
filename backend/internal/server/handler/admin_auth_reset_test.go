package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAdminAuthResetHandler(t *testing.T) {
	t.Run("RequestReset_Success", func(t *testing.T) {
		mockReqUC := &mock.MockAdminRequestPasswordResetUseCase{
			ExecuteFunc: func(ctx context.Context, email string) error {
				return nil
			},
		}
		mockReg := &mock.MockRegistry{RequestAdminPasswordResetUC: mockReqUC}
		h := handler.NewAdminAuthResetHandler(mockReg)

		reqBody := map[string]string{"email": "admin@example.com"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/admin/password-reset/request", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.RequestReset(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("VerifyToken_Success", func(t *testing.T) {
		mockReg := &mock.MockRegistry{} // VerifyToken is currently a placeholder
		h := handler.NewAdminAuthResetHandler(mockReg)

		// Body might not be parsed in VerifyToken current implementation (empty logic), but good to send valid req.
		reqBody := map[string]string{"token": "token123"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/admin/password-reset/verify", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.VerifyToken(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("ConfirmReset_Success", func(t *testing.T) {
		mockResetUC := &mock.MockAdminResetPasswordUseCase{
			ExecuteFunc: func(ctx context.Context, token, newPassword string) error {
				return nil
			},
		}
		mockReg := &mock.MockRegistry{ResetAdminPasswordUC: mockResetUC}
		h := handler.NewAdminAuthResetHandler(mockReg)

		reqBody := map[string]string{"token": "token123", "new_password": "newpass"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/admin/password-reset/confirm", bytes.NewReader(body))
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

		req := httptest.NewRequest(http.MethodGet, "/api/admin/password-reset/request", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}
