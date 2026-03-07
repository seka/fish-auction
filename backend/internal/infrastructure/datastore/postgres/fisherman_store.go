package postgres

import (
	"context"
	"database/sql"
	"errors"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type fishermanStore struct {
	db datastore.Database
}

func NewFishermanStore(db datastore.Database) repository.FishermanRepository {
	return &fishermanStore{
		db: db,
	}
}

func (r *fishermanStore) Create(ctx context.Context, name string) (*model.Fisherman, error) {
	e := entity.Fisherman{Name: name}
	if err := e.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRow(ctx, "INSERT INTO fishermen (name) VALUES ($1) RETURNING id, name", name).Scan(&e.ID, &e.Name)
	if err != nil {
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *fishermanStore) List(ctx context.Context) ([]model.Fisherman, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name FROM fishermen WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var fishermen []model.Fisherman
	for rows.Next() {
		var e entity.Fisherman
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, err
		}
		fishermen = append(fishermen, *e.ToModel())
	}
	return fishermen, rows.Err()
}

func (r *fishermanStore) FindByID(ctx context.Context, id int) (*model.Fisherman, error) {
	var e entity.Fisherman
	err := r.db.QueryRow(ctx,
		"SELECT id, name FROM fishermen WHERE id = $1",
		id,
	).Scan(&e.ID, &e.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Fisherman", ID: id}
		}
		return nil, err
	}

	return e.ToModel(), nil
}

func (r *fishermanStore) Delete(ctx context.Context, id int) error {
	_, err := r.db.Execute(ctx, "UPDATE fishermen SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", id)
	return err
}
