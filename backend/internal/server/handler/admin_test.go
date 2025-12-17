package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAdminHandler_UpdatePassword(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUpdateUC := &mock.MockAdminUpdatePasswordUseCase{
			ExecuteFunc: func(ctx context.Context, id int, currentPassword, newPassword string) error {
				if id != 1 {
					t.Errorf("expected id 1, got %d", id)
				}
				return nil
			},
		}
		mockReg := &mock.MockRegistry{UpdateAdminPasswordUC: mockUpdateUC}
		h := handler.NewAdminHandler(mockReg)

		reqBody := dto.UpdatePasswordRequest{
			CurrentPassword: "old",
			NewPassword:     "new",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/api/admin/password", bytes.NewReader(body))
		req.AddCookie(&http.Cookie{Name: "admin_session", Value: "authenticated"})
		req.AddCookie(&http.Cookie{Name: "admin_id", Value: "1"})
		w := httptest.NewRecorder()

		h.UpdatePassword(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := handler.NewAdminHandler(mockReg)

		req := httptest.NewRequest(http.MethodPut, "/api/admin/password", nil)
		w := httptest.NewRecorder()

		h.UpdatePassword(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}

func TestAdminHandler_RegisterRoutes(t *testing.T) {
	t.Run("MethodNotAllowed", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := handler.NewAdminHandler(mockReg)
		mux := http.NewServeMux()
		h.RegisterRoutes(mux)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/password", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status 405, got %d", w.Code)
		}
	})
}
