package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/entity"
	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestAuthHandler_Login(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockLoginUC := &mock.MockLoginUseCase{
			ExecuteFunc: func(ctx context.Context, email, password string) (*entity.Admin, error) {
				return &entity.Admin{ID: 1, Email: email}, nil
			},
		}
		mockReg := &mock.MockRegistry{LoginUC: mockLoginUC}
		h := handler.NewAuthHandler(mockReg)

		reqBody := dto.LoginRequest{Email: "admin@example.com", Password: "password"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Login(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		// Check cookies
		cookies := w.Result().Cookies()
		foundSession := false
		foundID := false
		for _, c := range cookies {
			if c.Name == "admin_session" && c.Value == "authenticated" {
				foundSession = true
			}
			if c.Name == "admin_id" && c.Value == "1" {
				foundID = true
			}
		}
		if !foundSession {
			t.Error("expected admin_session cookie")
		}
		if !foundID {
			t.Error("expected admin_id cookie")
		}
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mockLoginUC := &mock.MockLoginUseCase{
			ExecuteFunc: func(ctx context.Context, email, password string) (*entity.Admin, error) {
				return nil, nil // Returns nil, nil for invalid credentials (as per implementation inspection)
			},
		}
		mockReg := &mock.MockRegistry{LoginUC: mockLoginUC}
		h := handler.NewAuthHandler(mockReg)

		reqBody := dto.LoginRequest{Email: "admin@example.com", Password: "wrong"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Login(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}
