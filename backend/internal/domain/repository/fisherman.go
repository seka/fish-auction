package repository

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
)

type FishermanRepository interface {
	Create(ctx context.Context, name string) (*model.Fisherman, error)
	List(ctx context.Context) ([]model.Fisherman, error)
}
