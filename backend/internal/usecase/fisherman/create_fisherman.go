package fisherman

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateFishermanUseCase handles the creation of fishermen
type CreateFishermanUseCase struct {
	repo repository.FishermanRepository
}

// NewCreateFishermanUseCase creates a new instance of CreateFishermanUseCase
func NewCreateFishermanUseCase(repo repository.FishermanRepository) *CreateFishermanUseCase {
	return &CreateFishermanUseCase{repo: repo}
}

// Execute creates a new fisherman
func (uc *CreateFishermanUseCase) Execute(ctx context.Context, name string) (*model.Fisherman, error) {
	return uc.repo.Create(ctx, name)
}
