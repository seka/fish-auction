package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

type adminRepository struct {
	db datastore.Database
}

// NewAdminRepository creates a new instance of AdminRepository
func NewAdminRepository(db datastore.Database) repository.AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) FindOneByEmail(ctx context.Context, email string) (*model.Admin, error) {
	query := `SELECT id, email, password_hash, created_at FROM admins WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	admin := &model.Admin{}
	err := row.Scan(&admin.ID, &admin.Email, &admin.PasswordHash, &admin.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Or specific ErrNotFound
		}
		return nil, err
	}
	return admin, nil
}

func (r *adminRepository) FindByID(ctx context.Context, id int) (*model.Admin, error) {
	query := `SELECT id, email, password_hash, created_at FROM admins WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	admin := &model.Admin{}
	err := row.Scan(&admin.ID, &admin.Email, &admin.PasswordHash, &admin.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Or specific ErrNotFound
		}
		return nil, err
	}
	return admin, nil
}

func (r *adminRepository) Create(ctx context.Context, admin *model.Admin) error {
	query := `INSERT INTO admins (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, admin.Email, admin.PasswordHash).Scan(&admin.ID, &admin.CreatedAt)
	return err
}

func (r *adminRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM admins`
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *adminRepository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	query := `UPDATE admins SET password_hash = $1 WHERE id = $2`
	_, err := r.db.Execute(ctx, query, passwordHash, id)
	return err
}
