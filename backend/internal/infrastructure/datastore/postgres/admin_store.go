// Package postgres provides repository implementations using PostgreSQL.
package postgres

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
)

// AdminStore implements repository.AdminRepository using PostgreSQL.
type AdminStore struct {
	db datastore.Database
}

var _ repository.AdminRepository = (*AdminStore)(nil)

// NewAdminStore creates a new instance of AdminRepository
func NewAdminStore(db datastore.Database) *AdminStore {
	return &AdminStore{db: db}
}

// FindOneByEmail returns an admin by its email.
func (r *AdminStore) FindOneByEmail(ctx context.Context, email string) (*model.Admin, error) {
	query := `SELECT id, email, password_hash, created_at FROM admins WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	admin := &model.Admin{}
	err := row.Scan(&admin.ID, &admin.Email, &admin.PasswordHash, &admin.CreatedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Admin", 0, "FindOneByEmail")
	}
	return admin, nil
}

// FindByID returns an admin by its ID.
func (r *AdminStore) FindByID(ctx context.Context, id int) (*model.Admin, error) {
	query := `SELECT id, email, password_hash, created_at FROM admins WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	admin := &model.Admin{}
	err := row.Scan(&admin.ID, &admin.Email, &admin.PasswordHash, &admin.CreatedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Admin", id, "FindByID")
	}
	return admin, nil
}

// Create stores a new admin.
func (r *AdminStore) Create(ctx context.Context, admin *model.Admin) error {
	query := `INSERT INTO admins (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, admin.Email, admin.PasswordHash).Scan(&admin.ID, &admin.CreatedAt)
	if err != nil {
		return dserrors.HandleError(err, "Admin", 0, "Create")
	}
	return nil
}

// Count returns the total number of admins.
func (r *AdminStore) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM admins`
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, dserrors.HandleError(err, "Admin", 0, "Count")
	}
	return count, nil
}

// UpdatePassword updates the password hash of an admin.
func (r *AdminStore) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	query := `UPDATE admins SET password_hash = $1 WHERE id = $2`
	_, err := r.db.Execute(ctx, query, passwordHash, id)
	if err != nil {
		return dserrors.HandleError(err, "Admin", id, "UpdatePassword")
	}
	return nil
}
