package usecase

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// BuyerUseCase defines the interface for buyer-related business logic
type BuyerUseCase interface {
	Create(ctx context.Context, name string) (*model.Buyer, error)
	List(ctx context.Context) ([]model.Buyer, error)
}

type buyerInteractor struct {
	repo repository.BuyerRepository
}

func NewBuyerInteractor(repo repository.BuyerRepository) BuyerUseCase {
	return &buyerInteractor{repo: repo}
}

func (i *buyerInteractor) Create(ctx context.Context, name string) (*model.Buyer, error) {
	return i.repo.Create(ctx, name)
}

func (i *buyerInteractor) List(ctx context.Context) ([]model.Buyer, error) {
	return i.repo.List(ctx)
}
