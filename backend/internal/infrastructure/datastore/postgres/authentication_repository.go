package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type authenticationRepository struct {
	db *sql.DB
}

func NewAuthenticationRepository(db *sql.DB) repository.AuthenticationRepository {
	return &authenticationRepository{
		db: db,
	}
}

func (r *authenticationRepository) Create(ctx context.Context, auth *model.Authentication) (*model.Authentication, error) {
	e := entity.Authentication{
		BuyerID:      auth.BuyerID,
		Email:        auth.Email,
		PasswordHash: auth.PasswordHash,
		AuthType:     auth.AuthType,
	}
	if e.AuthType == "" {
		e.AuthType = "password"
	}

	if err := e.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRowContext(ctx,
		`INSERT INTO authentications (buyer_id, email, password_hash, auth_type, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
		 RETURNING id, created_at, updated_at`,
		e.BuyerID, e.Email, e.PasswordHash, e.AuthType).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *authenticationRepository) FindByEmail(ctx context.Context, email string) (*model.Authentication, error) {
	var e entity.Authentication
	err := r.db.QueryRowContext(ctx,
		`SELECT id, buyer_id, email, password_hash, auth_type, failed_attempts, locked_until, last_login_at, created_at, updated_at 
		 FROM authentications WHERE email = $1`,
		email,
	).Scan(&e.ID, &e.BuyerID, &e.Email, &e.PasswordHash, &e.AuthType, &e.FailedAttempts, &e.LockedUntil, &e.LastLoginAt, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Authentication", ID: 0}
		}
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *authenticationRepository) FindByBuyerID(ctx context.Context, buyerID int) (*model.Authentication, error) {
	var e entity.Authentication
	err := r.db.QueryRowContext(ctx,
		`SELECT id, buyer_id, email, password_hash, auth_type, failed_attempts, locked_until, last_login_at, created_at, updated_at 
		 FROM authentications WHERE buyer_id = $1`,
		buyerID,
	).Scan(&e.ID, &e.BuyerID, &e.Email, &e.PasswordHash, &e.AuthType, &e.FailedAttempts, &e.LockedUntil, &e.LastLoginAt, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Authentication", ID: buyerID}
		}
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *authenticationRepository) UpdateLoginSuccess(ctx context.Context, id int, loginAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE authentications 
		 SET last_login_at = $1, failed_attempts = 0, locked_until = NULL, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $2`,
		loginAt, id)
	return err
}

func (r *authenticationRepository) IncrementFailedAttempts(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE authentications 
		 SET failed_attempts = failed_attempts + 1, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $1`,
		id)
	return err
}

func (r *authenticationRepository) ResetFailedAttempts(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE authentications 
		 SET failed_attempts = 0, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $1`,
		id)
	return err
}

func (r *authenticationRepository) LockAccount(ctx context.Context, id int, until time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE authentications 
		 SET locked_until = $1, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $2`,
		until, id)
	return err
}

func (r *authenticationRepository) UpdatePassword(ctx context.Context, buyerID int, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE authentications 
		 SET password_hash = $1, updated_at = CURRENT_TIMESTAMP 
		 WHERE buyer_id = $2`,
		passwordHash, buyerID)
	return err
}
