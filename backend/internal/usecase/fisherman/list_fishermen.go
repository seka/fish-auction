package fisherman

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// ListFishermenUseCase defines the interface for listing fishermen
type ListFishermenUseCase interface {
	Execute(ctx context.Context) ([]model.Fisherman, error)
}

// listFishermenUseCase handles listing fishermen
type listFishermenUseCase struct {
	repo repository.FishermanRepository
}

// NewListFishermenUseCase creates a new instance of ListFishermenUseCase
func NewListFishermenUseCase(repo repository.FishermanRepository) ListFishermenUseCase {
	return &listFishermenUseCase{repo: repo}
}

// Execute lists all fishermen
func (uc *listFishermenUseCase) Execute(ctx context.Context) ([]model.Fisherman, error) {
	return uc.repo.List(ctx)
}
