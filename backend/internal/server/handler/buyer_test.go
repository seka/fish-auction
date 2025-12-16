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
	t.Run("Success", func(t *testing.T) {
		mockCreateUC := &mock.MockCreateBuyerUseCase{
			ExecuteFunc: func(ctx context.Context, name, email, password, organization, contactInfo string) (*model.Buyer, error) {
				return &model.Buyer{ID: 1, Name: name, Organization: organization}, nil
			},
		}
		mockReg := &mock.MockRegistry{CreateBuyerUC: mockCreateUC}
		h := handler.NewBuyerHandler(mockReg)

		reqBody := dto.CreateBuyerRequest{
			Name:         "Buyer 1",
			Email:        "buyer@example.com",
			Password:     "password",
			Organization: "Org 1",
			ContactInfo:  "Contact 1",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/buyers", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Create(w, req)

		if w.Code != http.StatusOK { // Handler returns 200 by default usually, check logic
			// handler/buyer.go Create uses json.NewEncoder(w).Encode(resp), so 200.
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestBuyerHandler_Login(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockLoginUC := &mock.MockLoginBuyerUseCase{
			ExecuteFunc: func(ctx context.Context, email, password string) (*model.Buyer, error) {
				return &model.Buyer{ID: 1, Name: "Buyer 1"}, nil
			},
		}
		mockReg := &mock.MockRegistry{LoginBuyerUC: mockLoginUC}
		h := handler.NewBuyerHandler(mockReg)

		reqBody := dto.LoginBuyerRequest{Email: "buyer@example.com", Password: "password"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/buyers/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.Login(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		cookies := w.Result().Cookies()
		foundSession := false
		foundID := false
		for _, c := range cookies {
			if c.Name == "buyer_session" && c.Value == "authenticated" {
				foundSession = true
			}
			if c.Name == "buyer_id" && c.Value == "1" {
				foundID = true
			}
		}
		if !foundSession {
			t.Error("expected buyer_session cookie")
		}
		if !foundID {
			t.Error("expected buyer_id cookie")
		}
	})
}

func TestBuyerHandler_GetMyPurchases(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockGetPurchasesUC := &mock.MockGetBuyerPurchasesUseCase{
			ExecuteFunc: func(ctx context.Context, buyerID int) ([]model.Purchase, error) {
				if buyerID != 1 {
					return nil, errors.New("wrong buyer ID")
				}
				return []model.Purchase{
					{ID: 1, ItemID: 1, FishType: "Tuna", Quantity: 10, Unit: "kg", Price: 1000, BuyerID: 1, CreatedAt: time.Now()},
				}, nil
			},
		}
		mockReg := &mock.MockRegistry{GetBuyerPurchasesUC: mockGetPurchasesUC}
		h := handler.NewBuyerHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/buyers/me/purchases", nil)
		req.AddCookie(&http.Cookie{Name: "buyer_session", Value: "authenticated"})
		req.AddCookie(&http.Cookie{Name: "buyer_id", Value: "1"})
		w := httptest.NewRecorder()

		h.GetMyPurchases(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("Unauthorized_NoCookie", func(t *testing.T) {
		mockReg := &mock.MockRegistry{}
		h := handler.NewBuyerHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/buyers/me/purchases", nil)
		w := httptest.NewRecorder()

		h.GetMyPurchases(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}
