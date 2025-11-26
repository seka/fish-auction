package bid_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/usecase/bid"
	mock "github.com/seka/fish-auction/backend/internal/usecase/testing"
)

func TestCreateBidUseCase_Execute(t *testing.T) {
	updateStatusErr := errors.New("update status failed")
	createBidErr := errors.New("create bid failed")
	txErr := errors.New("tx error")

	tests := []struct {
		name              string
		input             *model.Bid
		updateStatusErr   error
		createErr         error
		txErr             error
		wantID            int
		wantErr           error
		wantUpdateCalled  bool
		wantCreateCalled  bool
		wantTxCalled      bool
	}{
		{
			name: "Success",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			wantID:           1,
			wantUpdateCalled: true,
			wantCreateCalled: true,
			wantTxCalled:     true,
		},
		{
			name: "Error_UpdateStatusFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			updateStatusErr:  updateStatusErr,
			wantErr:          updateStatusErr,
			wantUpdateCalled: true,
			wantTxCalled:     true,
		},
		{
			name: "Error_CreateBidFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			createErr:        createBidErr,
			wantErr:          createBidErr,
			wantUpdateCalled: true,
			wantCreateCalled: true,
			wantTxCalled:     true,
		},
		{
			name: "Error_TransactionManagerFails",
			input: &model.Bid{
				ItemID:  1,
				BuyerID: 1,
				Price:   1000,
			},
			txErr:      txErr,
			wantErr:    txErr,
			wantTxCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCalled := false
			createCalled := false
			txCalled := false

			mockItemRepo := &mock.MockItemRepository{
				UpdateStatusFunc: func(ctx context.Context, id int, status model.ItemStatus) error {
					updateCalled = true
					if status != model.ItemStatusSold {
						t.Fatalf("unexpected status passed: %v", status)
					}
					if id != tt.input.ItemID {
						t.Fatalf("unexpected item id: %d", id)
					}
					return tt.updateStatusErr
				},
			}

			mockBidRepo := &mock.MockBidRepository{
				CreateFunc: func(ctx context.Context, b *model.Bid) (*model.Bid, error) {
					createCalled = true
					if b != tt.input {
						t.Fatalf("bid pointer mismatch")
					}
					if tt.createErr != nil {
						return nil, tt.createErr
					}
					cloned := *b
					cloned.ID = tt.wantID
					return &cloned, nil
				},
			}

			mockTxMgr := &mock.MockTransactionManager{
				WithTransactionFunc: func(ctx context.Context, fn func(ctx context.Context) error) error {
					txCalled = true
					if tt.txErr != nil {
						return tt.txErr
					}
					return fn(ctx)
				},
			}

			uc := bid.NewCreateBidUseCase(mockItemRepo, mockBidRepo, mockTxMgr)
			created, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if created != nil {
					t.Fatalf("expected nil result, got %+v", created)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if created == nil || created.ID != tt.wantID {
					t.Fatalf("unexpected created bid %+v", created)
				}
			}

			if updateCalled != tt.wantUpdateCalled {
				t.Fatalf("UpdateStatus called = %v, want %v", updateCalled, tt.wantUpdateCalled)
			}
			if createCalled != tt.wantCreateCalled {
				t.Fatalf("Create called = %v, want %v", createCalled, tt.wantCreateCalled)
			}
			if txCalled != tt.wantTxCalled {
				t.Fatalf("WithTransaction called = %v, want %v", txCalled, tt.wantTxCalled)
			}
		})
	}
}
