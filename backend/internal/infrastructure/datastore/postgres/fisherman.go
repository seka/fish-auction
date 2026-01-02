package postgres

import (
	"context"
	"database/sql"
	"errors"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	cache "github.com/seka/fish-auction/backend/internal/infrastructure/cache/redis"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type fishermanRepository struct {
	db    *sql.DB
	cache cache.FishermanCache
}

func NewFishermanRepository(db *sql.DB, fishermanCache cache.FishermanCache) repository.FishermanRepository {
	return &fishermanRepository{
		db:    db,
		cache: fishermanCache,
	}
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

func (r *fishermanRepository) FindByID(ctx context.Context, id int) (*model.Fisherman, error) {
	// キャッシュを確認
	if fisherman, err := r.cache.Get(ctx, id); err == nil && fisherman != nil {
		return fisherman, nil
	}

	// DBから取得
	var e entity.Fisherman
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name FROM fishermen WHERE id = $1",
		id,
	).Scan(&e.ID, &e.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Fisherman", ID: id}
		}
		return nil, err
	}

	fisherman := e.ToModel()

	// キャッシュに保存（エラーは無視）
	_ = r.cache.Set(ctx, id, fisherman)

	return fisherman, nil
}
