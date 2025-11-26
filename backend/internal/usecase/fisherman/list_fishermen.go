package fisherman

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListFishermenUseCase handles listing fishermen
type ListFishermenUseCase struct {
	repo repository.FishermanRepository
}

// NewListFishermenUseCase creates a new instance of ListFishermenUseCase
func NewListFishermenUseCase(repo repository.FishermanRepository) *ListFishermenUseCase {
	return &ListFishermenUseCase{repo: repo}
}

// Execute lists all fishermen
func (uc *ListFishermenUseCase) Execute(ctx context.Context) ([]model.Fisherman, error) {
	return uc.repo.List(ctx)
}
