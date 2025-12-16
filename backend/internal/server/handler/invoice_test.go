package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	mock "github.com/seka/fish-auction/backend/internal/server/testing"
)

func TestInvoiceHandler_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockListUC := &mock.MockListInvoicesUseCase{
			ExecuteFunc: func(ctx context.Context) ([]model.InvoiceItem, error) {
				return []model.InvoiceItem{
					{BuyerID: 1, BuyerName: "B1", TotalAmount: 100},
				}, nil
			},
		}
		mockReg := &mock.MockRegistry{ListInvoicesUC: mockListUC}
		h := handler.NewInvoiceHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/invoices", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockListUC := &mock.MockListInvoicesUseCase{
			ExecuteFunc: func(ctx context.Context) ([]model.InvoiceItem, error) {
				return nil, errors.New("db error")
			},
		}
		mockReg := &mock.MockRegistry{ListInvoicesUC: mockListUC}
		h := handler.NewInvoiceHandler(mockReg)

		req := httptest.NewRequest(http.MethodGet, "/api/invoices", nil)
		w := httptest.NewRecorder()

		h.List(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}
