package postgres

import (
	"context"
	"database/sql"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type BuyerRepository struct {
	db *sql.DB
}

func NewBuyerRepository(db *sql.DB) repository.BuyerRepository {
	return &BuyerRepository{db: db}
}

func (r *BuyerRepository) Create(ctx context.Context, name string) (*model.Buyer, error) {
	e := entity.Buyer{Name: name}
	if err := e.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRowContext(ctx, "INSERT INTO buyers (name) VALUES ($1) RETURNING id, name", name).Scan(&e.ID, &e.Name)
	if err != nil {
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *BuyerRepository) List(ctx context.Context) ([]model.Buyer, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM buyers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buyers []model.Buyer
	for rows.Next() {
		var e entity.Buyer
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, err
		}
		buyers = append(buyers, *e.ToModel())
	}
	return buyers, nil
}
