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
	repo repository.BuyerRepository
}

func NewGetBuyerUseCase(repo repository.BuyerRepository) GetBuyerUseCase {
	return &getBuyerUseCase{repo: repo}
}

func (uc *getBuyerUseCase) Execute(ctx context.Context, id int) (*model.Buyer, error) {
	return uc.repo.FindByID(ctx, id)
}
