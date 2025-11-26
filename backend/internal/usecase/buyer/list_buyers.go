package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListBuyersUseCase handles listing buyers
type ListBuyersUseCase struct {
	repo repository.BuyerRepository
}

// NewListBuyersUseCase creates a new instance of ListBuyersUseCase
func NewListBuyersUseCase(repo repository.BuyerRepository) *ListBuyersUseCase {
	return &ListBuyersUseCase{repo: repo}
}

// Execute lists all buyers
func (uc *ListBuyersUseCase) Execute(ctx context.Context) ([]model.Buyer, error) {
	return uc.repo.List(ctx)
}
