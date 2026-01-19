package item

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type DeleteItemUseCase interface {
	Execute(ctx context.Context, id int) error
}

type deleteItemUseCase struct {
	repo repository.ItemRepository
}

func NewDeleteItemUseCase(repo repository.ItemRepository) DeleteItemUseCase {
	return &deleteItemUseCase{repo: repo}
}

func (uc *deleteItemUseCase) Execute(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
