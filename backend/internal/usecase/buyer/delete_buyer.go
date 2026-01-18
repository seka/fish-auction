package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type DeleteBuyerUseCase interface {
	Execute(ctx context.Context, id int) error
}

type deleteBuyerUseCase struct {
	repo repository.BuyerRepository
}

func NewDeleteBuyerUseCase(repo repository.BuyerRepository) DeleteBuyerUseCase {
	return &deleteBuyerUseCase{repo: repo}
}

func (uc *deleteBuyerUseCase) Execute(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
