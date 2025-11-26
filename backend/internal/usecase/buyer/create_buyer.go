package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateBuyerUseCase handles the creation of buyers
type CreateBuyerUseCase struct {
	repo repository.BuyerRepository
}

// NewCreateBuyerUseCase creates a new instance of CreateBuyerUseCase
func NewCreateBuyerUseCase(repo repository.BuyerRepository) *CreateBuyerUseCase {
	return &CreateBuyerUseCase{repo: repo}
}

// Execute creates a new buyer
func (uc *CreateBuyerUseCase) Execute(ctx context.Context, name string) (*model.Buyer, error) {
	return uc.repo.Create(ctx, name)
}
