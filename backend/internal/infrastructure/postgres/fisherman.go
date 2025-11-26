package postgres

import (
	"context"
	"database/sql"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type fishermanRepository struct {
	db *sql.DB
}

func NewFishermanRepository(db *sql.DB) repository.FishermanRepository {
	return &fishermanRepository{db: db}
}

func (r *fishermanRepository) Create(ctx context.Context, name string) (*model.Fisherman, error) {
	e := entity.Fisherman{Name: name}
	if err := e.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRowContext(ctx, "INSERT INTO fishermen (name) VALUES ($1) RETURNING id, name", name).Scan(&e.ID, &e.Name)
	if err != nil {
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *fishermanRepository) List(ctx context.Context) ([]model.Fisherman, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM fishermen")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fishermen []model.Fisherman
	for rows.Next() {
		var e entity.Fisherman
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, err
		}
		fishermen = append(fishermen, *e.ToModel())
	}
	return fishermen, nil
}
