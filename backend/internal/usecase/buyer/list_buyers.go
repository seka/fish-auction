package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListBuyersUseCase defines the interface for listing buyers.
type ListBuyersUseCase interface {
	// Execute lists all buyers.
	Execute(ctx context.Context) ([]model.Buyer, error)
}

// ListBuyersUseCase handles listing buyers
type listBuyersUseCase struct {
	repo repository.BuyerRepository
}

var _ ListBuyersUseCase = (*listBuyersUseCase)(nil)

// NewListBuyersUseCase creates a new instance of ListBuyersUseCase
func NewListBuyersUseCase(buyerRepo repository.BuyerRepository) ListBuyersUseCase {
	return &listBuyersUseCase{repo: buyerRepo}
}

// Execute lists all buyers
func (uc *listBuyersUseCase) Execute(ctx context.Context) ([]model.Buyer, error) {
	return uc.repo.List(ctx)
}
