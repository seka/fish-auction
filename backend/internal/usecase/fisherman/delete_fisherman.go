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

var _ DeleteFishermanUseCase = (*deleteFishermanUseCase)(nil)

type deleteFishermanUseCase struct {
	repo repository.FishermanRepository
}

// NewDeleteFishermanUseCase creates a new instance of DeleteFishermanUseCase
func NewDeleteFishermanUseCase(repo repository.FishermanRepository) *deleteFishermanUseCase {
	return &deleteFishermanUseCase{repo: repo}
}

func (uc *deleteFishermanUseCase) Execute(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}
