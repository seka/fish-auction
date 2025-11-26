package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateBuyerUseCase defines the interface for creating buyers
type CreateBuyerUseCase interface {
	Execute(ctx context.Context, name string) (*model.Buyer, error)
}

// createBuyerUseCase handles the creation of buyers
type createBuyerUseCase struct {
	repo repository.BuyerRepository
}

// NewCreateBuyerUseCase creates a new instance of CreateBuyerUseCase
func NewCreateBuyerUseCase(repo repository.BuyerRepository) CreateBuyerUseCase {
	return &createBuyerUseCase{repo: repo}
}

// Execute creates a new buyer
func (uc *createBuyerUseCase) Execute(ctx context.Context, name string) (*model.Buyer, error) {
	return uc.repo.Create(ctx, name)
}
