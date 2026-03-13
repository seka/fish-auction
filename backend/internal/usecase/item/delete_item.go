package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type DeleteItemUseCase interface {
	Execute(ctx context.Context, id int) error
}

var _ DeleteItemUseCase = (*deleteItemUseCase)(nil)

type deleteItemUseCase struct {
	repo repository.ItemRepository
}

// NewDeleteItemUseCase creates a new instance of DeleteItemUseCase
func NewDeleteItemUseCase(repo repository.ItemRepository) *deleteItemUseCase {
	return &deleteItemUseCase{repo: repo}
}

func (uc *deleteItemUseCase) Execute(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
