package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListBuyersUseCase defines the interface for listing buyers
type ListBuyersUseCase interface {
	Execute(ctx context.Context) ([]model.Buyer, error)
}

// listBuyersUseCase handles listing buyers
type listBuyersUseCase struct {
	repo repository.BuyerRepository
}

// NewListBuyersUseCase creates a new instance of ListBuyersUseCase
func NewListBuyersUseCase(repo repository.BuyerRepository) ListBuyersUseCase {
	return &listBuyersUseCase{repo: repo}
}

// Execute lists all buyers
func (uc *listBuyersUseCase) Execute(ctx context.Context) ([]model.Buyer, error) {
	return uc.repo.List(ctx)
}
