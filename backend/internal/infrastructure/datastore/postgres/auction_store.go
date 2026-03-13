package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/errors"
)

var _ repository.AuctionRepository = (*auctionStore)(nil)

type auctionStore struct {
	db datastore.Database
}

// NewAuctionStore creates a new instance of AuctionRepository
func NewAuctionStore(db datastore.Database) *auctionStore {
	return &auctionStore{db: db}
}

func (r *auctionStore) Create(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
	query := `INSERT INTO auctions (venue_id, auction_date, start_time, end_time, status)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at`

	var a model.Auction
	err := r.db.QueryRow(ctx, query,
		auction.VenueID, auction.AuctionDate, auction.StartTime, auction.EndTime, auction.Status).
		Scan(&a.ID, &a.VenueID, &a.AuctionDate, &a.StartTime, &a.EndTime, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == dserrors.ErrCodeUniqueViolation {
			return nil, &apperrors.ConflictError{Message: fmt.Sprintf("Auction already exists for venue %d on %s", auction.VenueID, auction.AuctionDate.Format("2006-01-02"))}
		}
		return nil, err
	}
	return &a, nil
}

func (r *auctionStore) GetByID(ctx context.Context, id int) (*model.Auction, error) {
	query := `SELECT id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at
			  FROM auctions WHERE id = $1`

	var a model.Auction
	err := r.db.QueryRow(ctx, query, id).
		Scan(&a.ID, &a.VenueID, &a.AuctionDate, &a.StartTime, &a.EndTime, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.NotFoundError{Resource: "Auction", ID: id}
		}
		return nil, err
	}
	return &a, nil
}

func (r *auctionStore) List(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	query := `SELECT id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at
			  FROM auctions`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters != nil {
		if filters.VenueID != nil {
			conditions = append(conditions, fmt.Sprintf("venue_id = $%d", argIndex))
			args = append(args, *filters.VenueID)
			argIndex++
		}
		if filters.AuctionDate != nil {
			conditions = append(conditions, fmt.Sprintf("auction_date = $%d", argIndex))
			args = append(args, *filters.AuctionDate)
			argIndex++
		}
		if filters.Status != nil {
			conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
			args = append(args, *filters.Status)
			argIndex++
		}
		if filters.StartDate != nil {
			conditions = append(conditions, fmt.Sprintf("auction_date >= $%d", argIndex))
			args = append(args, *filters.StartDate)
			argIndex++
		}
		if filters.EndDate != nil {
			conditions = append(conditions, fmt.Sprintf("auction_date <= $%d", argIndex))
			args = append(args, *filters.EndDate)
			argIndex++
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY auction_date DESC, created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var auctions []model.Auction
	for rows.Next() {
		var a model.Auction
		if err := rows.Scan(&a.ID, &a.VenueID, &a.AuctionDate, &a.StartTime, &a.EndTime, &a.Status, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		auctions = append(auctions, a)
	}
	return auctions, rows.Err()
}

func (r *auctionStore) ListByVenue(ctx context.Context, venueID int) ([]model.Auction, error) {
	filters := &repository.AuctionFilters{
		VenueID: &venueID,
	}
	return r.List(ctx, filters)
}

func (r *auctionStore) Update(ctx context.Context, auction *model.Auction) error {
	query := `UPDATE auctions
			  SET venue_id = $1, auction_date = $2, start_time = $3, end_time = $4, status = $5, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $6`

	rowsAffected, err := r.db.Execute(ctx, query,
		auction.VenueID, auction.AuctionDate, auction.StartTime, auction.EndTime, auction.Status, auction.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == dserrors.ErrCodeUniqueViolation {
			return &apperrors.ConflictError{Message: "Auction already exists for this venue and time"}
		}
		return err
	}

	if rowsAffected == 0 {
		return &apperrors.NotFoundError{Resource: "Auction", ID: auction.ID}
	}
	return nil
}

func (r *auctionStore) UpdateStatus(ctx context.Context, id int, status model.AuctionStatus) error {
	query := `UPDATE auctions SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	rowsAffected, err := r.db.Execute(ctx, query, status, id)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &apperrors.NotFoundError{Resource: "Auction", ID: id}
	}
	return nil
}

// Delete はセリをデータベースから削除します。
// CASCADE削除: この操作により以下のデータも自動的に削除されます:
//   - このセリに関連付けられたすべての出品
//
// 注意: 出品に入札（transactions）が存在する場合、入札履歴を保護するため削除は失敗します。
func (r *auctionStore) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM auctions WHERE id = $1`
	_, err := r.db.Execute(ctx, query, id)
	return err
}
