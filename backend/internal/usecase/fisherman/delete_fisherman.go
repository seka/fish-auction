package fisherman

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// DeleteFishermanUseCase defines the interface for deleting a fisherman.
type DeleteFishermanUseCase interface {
	// Execute deletes a fisherman by ID.
	Execute(ctx context.Context, id int) error
}

type deleteFishermanUseCase struct {
	repo repository.FishermanRepository
}

func NewDeleteFishermanUseCase(repo repository.FishermanRepository) DeleteFishermanUseCase {
	return &deleteFishermanUseCase{repo: repo}
}

func (uc *deleteFishermanUseCase) Execute(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
