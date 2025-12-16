package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type BuyerRepository interface {
	Create(ctx context.Context, buyer *model.Buyer) (*model.Buyer, error)
	List(ctx context.Context) ([]model.Buyer, error)
	FindByID(ctx context.Context, id int) (*model.Buyer, error)
	FindByName(ctx context.Context, name string) (*model.Buyer, error)
	FindByEmail(ctx context.Context, email string) (*model.Buyer, error)
}
