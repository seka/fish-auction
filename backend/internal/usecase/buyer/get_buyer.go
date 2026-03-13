package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type GetBuyerUseCase interface {
	Execute(ctx context.Context, id int) (*model.Buyer, error)
}

type getBuyerUseCase struct {
	buyerRepo repository.BuyerRepository
}

var _ GetBuyerUseCase = (*getBuyerUseCase)(nil)

// NewGetBuyerUseCase creates a new instance of GetBuyerUseCase
func NewGetBuyerUseCase(buyerRepo repository.BuyerRepository) *getBuyerUseCase {
	return &getBuyerUseCase{buyerRepo: buyerRepo}
}

func (uc *getBuyerUseCase) Execute(ctx context.Context, id int) (*model.Buyer, error) {
	return uc.buyerRepo.FindByID(ctx, id)
}
