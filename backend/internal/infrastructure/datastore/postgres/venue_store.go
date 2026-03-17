package postgres

import (
	"context"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
)

var _ repository.VenueRepository = (*VenueStore)(nil)

// VenueStore implements repository.VenueRepository using PostgreSQL.
type VenueStore struct {
	db datastore.Database
}

// NewVenueStore creates a new instance of VenueRepository
func NewVenueStore(db datastore.Database) *VenueStore {
	return &VenueStore{db: db}
}

// Create stores a new venue.
func (r *VenueStore) Create(ctx context.Context, venue *model.Venue) (*model.Venue, error) {
	query := `INSERT INTO venues (name, location, description) VALUES ($1, $2, $3)
			  RETURNING id, name, location, description, created_at`

	var v model.Venue
	err := r.db.QueryRow(ctx, query, venue.Name, venue.Location, venue.Description).
		Scan(&v.ID, &v.Name, &v.Location, &v.Description, &v.CreatedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Venue", 0, "Create")
	}
	return &v, nil
}

// FindByID returns a venue by its ID.
func (r *VenueStore) FindByID(ctx context.Context, id int) (*model.Venue, error) {
	query := `SELECT id, name, location, description, created_at, deleted_at FROM venues WHERE id = $1`

	var v model.Venue
	err := r.db.QueryRow(ctx, query, id).
		Scan(&v.ID, &v.Name, &v.Location, &v.Description, &v.CreatedAt, &v.DeletedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Venue", id, "FindByID")
	}
	return &v, nil
}

// List returns all active venues.
func (r *VenueStore) List(ctx context.Context) ([]model.Venue, error) {
	query := `SELECT id, name, location, description, created_at, deleted_at FROM venues WHERE deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, dserrors.HandleError(err, "Venue", 0, "List")
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
	if err := rows.Err(); err != nil {
		return nil, dserrors.HandleError(err, "Venue", 0, "List")
	}
	return venues, nil
}

// Update updates an existing venue.
func (r *VenueStore) Update(ctx context.Context, venue *model.Venue) error {
	query := `UPDATE venues SET name = $1, location = $2, description = $3 WHERE id = $4 AND deleted_at IS NULL`

	rowsAffected, err := r.db.Execute(ctx, query, venue.Name, venue.Location, venue.Description, venue.ID)
	if err != nil {
		return dserrors.HandleError(err, "Venue", venue.ID, "Update")
	}

	if rowsAffected == 0 {
		return &apperrors.NotFoundError{Resource: "Venue", ID: venue.ID}
	}
	return nil
}

// Delete は会場を論理削除します。
// Delete marks a venue as deleted.
func (r *VenueStore) Delete(ctx context.Context, id int) error {
	query := `UPDATE venues SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	rowsAffected, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return dserrors.HandleError(err, "Venue", id, "Delete")
	}

	if rowsAffected == 0 {
		return &apperrors.NotFoundError{Resource: "Venue", ID: id}
	}
	return nil
}
