package usecase

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// FishermanUseCase defines the interface for fisherman-related business logic
type FishermanUseCase interface {
	Create(ctx context.Context, name string) (*model.Fisherman, error)
	List(ctx context.Context) ([]model.Fisherman, error)
}

type fishermanInteractor struct {
	repo repository.FishermanRepository
}

func NewFishermanInteractor(repo repository.FishermanRepository) FishermanUseCase {
	return &fishermanInteractor{repo: repo}
}

func (i *fishermanInteractor) Create(ctx context.Context, name string) (*model.Fisherman, error) {
	return i.repo.Create(ctx, name)
}

func (i *fishermanInteractor) List(ctx context.Context) ([]model.Fisherman, error) {
	return i.repo.List(ctx)
}
