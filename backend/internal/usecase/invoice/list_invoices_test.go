package invoice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/invoice"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestListInvoicesUseCase_Execute(t *testing.T) {
	tests := []struct {
		name     string
		invoices []model.InvoiceItem
		wantErr  error
	}{
		{
			name: "Success",
			invoices: []model.InvoiceItem{
				{BuyerID: 1, BuyerName: "Buyer1", TotalAmount: 1000},
				{BuyerID: 2, BuyerName: "Buyer2", TotalAmount: 2000},
			},
		},
		{
			name:    "Error",
			wantErr: errors.New("list invoices failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mock.MockBidRepository{
				ListInvoicesFunc: func(ctx context.Context) ([]model.InvoiceItem, error) {
					if tt.wantErr != nil {
						return nil, tt.wantErr
					}
					return tt.invoices, nil
				},
			}

			uc := invoice.NewListInvoicesUseCase(repo)
			got, err := uc.Execute(context.Background())

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(got) != len(tt.invoices) {
				t.Fatalf("expected %d invoices, got %d", len(tt.invoices), len(got))
			}
		})
	}
}
