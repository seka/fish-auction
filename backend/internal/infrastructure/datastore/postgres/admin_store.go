package postgres

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
)

type adminStore struct {
	db datastore.Database
}

var _ repository.AdminRepository = (*adminStore)(nil)

// NewAdminStore creates a new instance of AdminRepository.
func NewAdminStore(db datastore.Database) *adminStore {
	return &adminStore{db: db}
}

func (r *adminStore) FindOneByEmail(ctx context.Context, email string) (*model.Admin, error) {
	query := `SELECT id, email, password_hash, created_at FROM admins WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	admin := &model.Admin{}
	err := row.Scan(&admin.ID, &admin.Email, &admin.PasswordHash, &admin.CreatedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Admin", 0, "FindOneByEmail")
	}
	return admin, nil
}

func (r *adminStore) FindByID(ctx context.Context, id int) (*model.Admin, error) {
	query := `SELECT id, email, password_hash, created_at FROM admins WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	admin := &model.Admin{}
	err := row.Scan(&admin.ID, &admin.Email, &admin.PasswordHash, &admin.CreatedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Admin", id, "FindByID")
	}
	return admin, nil
}

func (r *adminStore) Create(ctx context.Context, admin *model.Admin) error {
	query := `INSERT INTO admins (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, admin.Email, admin.PasswordHash).Scan(&admin.ID, &admin.CreatedAt)
	if err != nil {
		return dserrors.HandleError(err, "Admin", 0, "Create")
	}
	return nil
}

func (r *adminStore) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM admins`
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, dserrors.HandleError(err, "Admin", 0, "Count")
	}
	return count, nil
}

func (r *adminStore) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	query := `UPDATE admins SET password_hash = $1 WHERE id = $2`
	_, err := r.db.Execute(ctx, query, passwordHash, id)
	if err != nil {
		return dserrors.HandleError(err, "Admin", id, "UpdatePassword")
	}
	return nil
}
