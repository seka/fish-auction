package postgres

import (
	"context"
	"database/sql"
	"errors"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/cache"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type buyerRepository struct {
	db    *sql.DB
	cache cache.BuyerCache
}

func NewBuyerRepository(db *sql.DB, buyerCache cache.BuyerCache) repository.BuyerRepository {
	return &buyerRepository{
		db:    db,
		cache: buyerCache,
	}
}

func (r *buyerRepository) Create(ctx context.Context, name string) (*model.Buyer, error) {
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

func (r *buyerRepository) List(ctx context.Context) ([]model.Buyer, error) {
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

func (r *buyerRepository) FindByID(ctx context.Context, id int) (*model.Buyer, error) {
	// キャッシュを確認
	if buyer, err := r.cache.Get(ctx, id); err == nil && buyer != nil {
		return buyer, nil
	}

	// DBから取得
	var e entity.Buyer
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name FROM buyers WHERE id = $1",
		id,
	).Scan(&e.ID, &e.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Buyer", ID: id}
		}
		return nil, err
	}

	buyer := e.ToModel()

	// キャッシュに保存（エラーは無視）
	_ = r.cache.Set(ctx, id, buyer)

	return buyer, nil
}
