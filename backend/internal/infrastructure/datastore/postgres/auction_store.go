package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
)

var _ repository.AuctionRepository = (*AuctionStore)(nil)

// AuctionStore implements repository.AuctionRepository using PostgreSQL.
type AuctionStore struct {
	db datastore.Database
}

// NewAuctionStore creates a new instance of AuctionRepository
func NewAuctionStore(db datastore.Database) *AuctionStore {
	return &AuctionStore{db: db}
}

// Create stores a new auction.
func (r *AuctionStore) Create(ctx context.Context, auction *model.Auction) (*model.Auction, error) {
	query := `INSERT INTO auctions (venue_id, auction_date, start_time, end_time, status)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at`

	var a model.Auction
	var auctionDate time.Time
	var startTime, endTime *time.Time
	err := r.db.QueryRow(ctx, query,
		auction.VenueID, auction.Period.AuctionDate, auction.Period.StartAt, auction.Period.EndAt, auction.Status).
		Scan(&a.ID, &a.VenueID, &auctionDate, &startTime, &endTime, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if dserrors.IsUniqueViolation(err) {
			return nil, &apperrors.ConflictError{Message: fmt.Sprintf("Auction already exists for venue %d on %s", auction.VenueID, auction.Period.AuctionDate.Format("2006-01-02"))}
		}
		return nil, dserrors.HandleError(err, "Auction", nil, "failed to create auction")
	}
	a.Period = model.NewAuctionPeriod(auctionDate, startTime, endTime)
	return &a, nil
}

// FindByID returns an auction by its ID.
func (r *AuctionStore) FindByID(ctx context.Context, id int) (*model.Auction, error) {
	query := `SELECT id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at
			  FROM auctions WHERE id = $1`

	var a model.Auction
	var auctionDate time.Time
	var startTime, endTime *time.Time
	err := r.db.QueryRow(ctx, query, id).
		Scan(&a.ID, &a.VenueID, &auctionDate, &startTime, &endTime, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Auction", id, "failed to get auction by ID")
	}
	a.Period = model.NewAuctionPeriod(auctionDate, startTime, endTime)
	return &a, nil
}

// FindByIDWithLock returns an auction by its ID with a lock.
func (r *AuctionStore) FindByIDWithLock(ctx context.Context, id int) (*model.Auction, error) {
	query := `SELECT id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at
			  FROM auctions WHERE id = $1 FOR UPDATE`

	var a model.Auction
	var auctionDate time.Time
	var startTime, endTime *time.Time
	err := r.db.QueryRow(ctx, query, id).
		Scan(&a.ID, &a.VenueID, &auctionDate, &startTime, &endTime, &a.Status, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Auction", id, "failed to get auction by ID with lock")
	}
	a.Period = model.NewAuctionPeriod(auctionDate, startTime, endTime)
	return &a, nil
}

// List returns a list of auctions based on the given filters.
func (r *AuctionStore) List(ctx context.Context, filters *repository.AuctionFilters) ([]model.Auction, error) {
	query := `SELECT id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at
			  FROM auctions`

	var conditions []string
	var args []any
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
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY auction_date DESC, created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, dserrors.HandleError(err, "Auction", nil, "failed to list auctions")
	}
	defer func() { _ = rows.Close() }()

	var auctions []model.Auction
	for rows.Next() {
		var a model.Auction
		var auctionDate time.Time
		var startTime, endTime *time.Time
		if err := rows.Scan(&a.ID, &a.VenueID, &auctionDate, &startTime, &endTime, &a.Status, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, dserrors.HandleError(err, "Auction", nil, "failed to scan auction row")
		}
		a.Period = model.NewAuctionPeriod(auctionDate, startTime, endTime)
		auctions = append(auctions, a)
	}
	if err := rows.Err(); err != nil {
		return nil, dserrors.HandleError(err, "Auction", nil, "failed to iterate auction rows")
	}
	return auctions, nil
}

// ListByVenue returns a list of auctions for the given venue ID.
func (r *AuctionStore) ListByVenue(ctx context.Context, venueID int) ([]model.Auction, error) {
	filters := &repository.AuctionFilters{
		VenueID: &venueID,
	}
	return r.List(ctx, filters)
}

// Update updates an existing auction.
func (r *AuctionStore) Update(ctx context.Context, auction *model.Auction) error {
	query := `UPDATE auctions
			  SET venue_id = $1, auction_date = $2, start_time = $3, end_time = $4, status = $5, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $6`

	rowsAffected, err := r.db.Execute(ctx, query,
		auction.VenueID, auction.Period.AuctionDate, auction.Period.StartAt, auction.Period.EndAt, auction.Status, auction.ID)
	if err != nil {
		if dserrors.IsUniqueViolation(err) {
			return &apperrors.ConflictError{Message: "Auction already exists for this venue and time"}
		}
		return dserrors.HandleError(err, "Auction", auction.ID, "failed to update auction")
	}

	if rowsAffected == 0 {
		return &apperrors.NotFoundError{Resource: "Auction", ID: auction.ID}
	}
	return nil
}

// UpdateStatus updates the status of an auction.
func (r *AuctionStore) UpdateStatus(ctx context.Context, id int, status model.AuctionStatus) error {
	query := `UPDATE auctions SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	rowsAffected, err := r.db.Execute(ctx, query, status, id)
	if err != nil {
		return dserrors.HandleError(err, "Auction", id, "failed to update auction status")
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
func (r *AuctionStore) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM auctions WHERE id = $1`
	_, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return dserrors.HandleError(err, "Auction", id, "failed to delete auction")
	}
	return nil
}
