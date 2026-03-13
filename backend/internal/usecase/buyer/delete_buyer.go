package buyer

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type DeleteBuyerUseCase interface {
	Execute(ctx context.Context, id int) error
}

var _ DeleteBuyerUseCase = (*deleteBuyerUseCase)(nil)

type deleteBuyerUseCase struct {
	repo repository.BuyerRepository
}

// NewDeleteBuyerUseCase creates a new instance of DeleteBuyerUseCase
func NewDeleteBuyerUseCase(repo repository.BuyerRepository) *deleteBuyerUseCase {
	return &deleteBuyerUseCase{repo: repo}
}

func (uc *deleteBuyerUseCase) Execute(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
