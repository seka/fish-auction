package postgres

import (
	"context"
	"database/sql"
	"errors"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

type venueStore struct {
	db datastore.Database
}

func NewVenueStore(db datastore.Database) repository.VenueRepository {
	return &venueStore{db: db}
}

func (r *venueStore) Create(ctx context.Context, venue *model.Venue) (*model.Venue, error) {
	query := `INSERT INTO venues (name, location, description) VALUES ($1, $2, $3)
			  RETURNING id, name, location, description, created_at`

	var v model.Venue
	err := r.db.QueryRow(ctx, query, venue.Name, venue.Location, venue.Description).
		Scan(&v.ID, &v.Name, &v.Location, &v.Description, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *venueStore) GetByID(ctx context.Context, id int) (*model.Venue, error) {
	query := `SELECT id, name, location, description, created_at, deleted_at FROM venues WHERE id = $1`

	var v model.Venue
	err := r.db.QueryRow(ctx, query, id).
		Scan(&v.ID, &v.Name, &v.Location, &v.Description, &v.CreatedAt, &v.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Venue", ID: id}
		}
		return nil, err
	}
	return &v, nil
}

func (r *venueStore) List(ctx context.Context) ([]model.Venue, error) {
	query := `SELECT id, name, location, description, created_at, deleted_at FROM venues WHERE deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var venues []model.Venue
	for rows.Next() {
		var v model.Venue
		if err := rows.Scan(&v.ID, &v.Name, &v.Location, &v.Description, &v.CreatedAt, &v.DeletedAt); err != nil {
			return nil, err
		}
		venues = append(venues, v)
	}
	return venues, rows.Err()
}

func (r *venueStore) Update(ctx context.Context, venue *model.Venue) error {
	query := `UPDATE venues SET name = $1, location = $2, description = $3 WHERE id = $4 AND deleted_at IS NULL`

	rowsAffected, err := r.db.Execute(ctx, query, venue.Name, venue.Location, venue.Description, venue.ID)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &apperrors.NotFoundError{Resource: "Venue", ID: venue.ID}
	}
	return nil
}

// Delete は会場を論理削除します。
func (r *venueStore) Delete(ctx context.Context, id int) error {
	query := `UPDATE venues SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	rowsAffected, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &apperrors.NotFoundError{Resource: "Venue", ID: id}
	}
	return nil
}
