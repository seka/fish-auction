package fisherman

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// CreateFishermanUseCase defines the interface for creating fishermen
type CreateFishermanUseCase interface {
	Execute(ctx context.Context, name string) (*model.Fisherman, error)
}

// createFishermanUseCase handles the creation of fishermen
type createFishermanUseCase struct {
	repo repository.FishermanRepository
}

var _ CreateFishermanUseCase = (*createFishermanUseCase)(nil)

// NewCreateFishermanUseCase creates a new instance of CreateFishermanUseCase
func NewCreateFishermanUseCase(repo repository.FishermanRepository) *createFishermanUseCase {
	return &createFishermanUseCase{repo: repo}
}

// Execute creates a new fisherman
func (uc *createFishermanUseCase) Execute(ctx context.Context, name string) (*model.Fisherman, error) {
	return uc.repo.Create(ctx, name)
}
