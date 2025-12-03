package postgres

import (
	"context"
	"database/sql"
	"errors"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type venueRepository struct {
	db *sql.DB
}

func NewVenueRepository(db *sql.DB) repository.VenueRepository {
	return &venueRepository{db: db}
}

func (r *venueRepository) Create(ctx context.Context, venue *model.Venue) (*model.Venue, error) {
	query := `INSERT INTO venues (name, location, description) VALUES ($1, $2, $3) 
			  RETURNING id, name, location, description, created_at`

	var v model.Venue
	err := r.db.QueryRowContext(ctx, query, venue.Name, venue.Location, venue.Description).
		Scan(&v.ID, &v.Name, &v.Location, &v.Description, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *venueRepository) GetByID(ctx context.Context, id int) (*model.Venue, error) {
	query := `SELECT id, name, location, description, created_at FROM venues WHERE id = $1`

	var v model.Venue
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&v.ID, &v.Name, &v.Location, &v.Description, &v.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Venue", ID: id}
		}
		return nil, err
	}
	return &v, nil
}

func (r *venueRepository) List(ctx context.Context) ([]model.Venue, error) {
	query := `SELECT id, name, location, description, created_at FROM venues ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var venues []model.Venue
	for rows.Next() {
		var v model.Venue
		if err := rows.Scan(&v.ID, &v.Name, &v.Location, &v.Description, &v.CreatedAt); err != nil {
			return nil, err
		}
		venues = append(venues, v)
	}
	return venues, nil
}

func (r *venueRepository) Update(ctx context.Context, venue *model.Venue) error {
	query := `UPDATE venues SET name = $1, location = $2, description = $3 WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query, venue.Name, venue.Location, venue.Description, venue.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &apperrors.NotFoundError{Resource: "Venue", ID: venue.ID}
	}
	return nil
}

// Delete は会場をデータベースから削除します。
// CASCADE削除: この操作により以下のデータも自動的に削除されます:
//   - この会場に関連付けられたすべてのセリ
//   - それらのセリに関連付けられたすべての出品
//
// 注意: 出品に入札（transactions）が存在する場合、入札履歴を保護するため削除は失敗します。
func (r *venueRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM venues WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &apperrors.NotFoundError{Resource: "Venue", ID: id}
	}
	return nil
}
