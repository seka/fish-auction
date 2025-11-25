package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type BuyerRepository interface {
	Create(ctx context.Context, name string) (*model.Buyer, error)
	List(ctx context.Context) ([]model.Buyer, error)
}
