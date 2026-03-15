package postgres

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

// BidStore implements repository.BidRepository using PostgreSQL.
type BidStore struct {
	db datastore.Database
}

// Ensure BidStore implements repository.BidRepository
var _ repository.BidRepository = (*BidStore)(nil)

// NewBidStore creates a new instance of BidRepository
func NewBidStore(db datastore.Database) *BidStore {
	return &BidStore{db: db}
}

// Create stores a new bid transaction.
func (r *BidStore) Create(ctx context.Context, bid *model.Bid) (*model.Bid, error) {

	e := entity.Bid{
		ItemID:  bid.ItemID,
		BuyerID: bid.BuyerID,
		Price:   bid.Price.Amount(),
	}
	if err := e.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRow(ctx,
		"INSERT INTO transactions (item_id, buyer_id, price) VALUES ($1, $2, $3) RETURNING id, item_id, buyer_id, price, created_at",
		bid.ItemID, bid.BuyerID, bid.Price.Amount(),
	).Scan(&e.ID, &e.ItemID, &e.BuyerID, &e.Price, &e.CreatedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Bid", 0, "Create")
	}
	return e.ToModel(), nil
}

// ListInvoices returns a list of invoice items based on bidding transactions.
func (r *BidStore) ListInvoices(ctx context.Context) ([]model.InvoiceItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT b.id, b.name, SUM(t.price) as total_price
		FROM transactions t
		JOIN buyers b ON t.buyer_id = b.id
		GROUP BY b.id, b.name
	`)
	if err != nil {
		return nil, dserrors.HandleError(err, "Invoice", 0, "ListInvoices")
	}
	defer func() { _ = rows.Close() }()

	var invoices []model.InvoiceItem
	for rows.Next() {
		var id int
		var name string
		var totalPrice int
		if err := rows.Scan(&id, &name, &totalPrice); err != nil {
			return nil, err
		}

		// 8% Tax
		totalAmount := int(float64(totalPrice) * 1.08)

		invoices = append(invoices, model.InvoiceItem{
			BuyerID:     id,
			BuyerName:   name,
			TotalAmount: totalAmount,
		})
	}
	return invoices, dserrors.HandleError(rows.Err(), "Invoice", 0, "ListInvoices")
}

// ListPurchasesByBuyerID returns all purchases for a specific buyer.
func (r *BidStore) ListPurchasesByBuyerID(ctx context.Context, buyerID int) ([]model.Purchase, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			t.id,
			t.item_id,
			ai.fish_type,
			ai.quantity,
			ai.unit,
			t.price,
			t.buyer_id,
			ai.auction_id,
			a.auction_date,
			t.created_at
		FROM transactions t
		JOIN auction_items ai ON t.item_id = ai.id
		JOIN auctions a ON ai.auction_id = a.id
		WHERE t.buyer_id = $1
		ORDER BY t.created_at DESC
	`, buyerID)
	if err != nil {
		return nil, dserrors.HandleError(err, "Purchase", buyerID, "ListPurchasesByBuyerID")
	}
	defer func() { _ = rows.Close() }()

	var purchases []model.Purchase
	for rows.Next() {
		var p model.Purchase
		if err := rows.Scan(
			&p.ID,
			&p.ItemID,
			&p.FishType,
			&p.Quantity,
			&p.Unit,
			&p.Price,
			&p.BuyerID,
			&p.AuctionID,
			&p.AuctionDate,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}
	return purchases, dserrors.HandleError(rows.Err(), "Purchase", buyerID, "ListPurchasesByBuyerID")
}

// ListAuctionsByBuyerID returns all auctions in which a specific buyer participated.
func (r *BidStore) ListAuctionsByBuyerID(ctx context.Context, buyerID int) ([]model.Auction, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT
			a.id,
			a.venue_id,
			a.auction_date,
			a.start_time,
			a.end_time,
			a.status,
			a.created_at,
			a.updated_at
		FROM auctions a
		JOIN auction_items ai ON a.id = ai.auction_id
		JOIN transactions t ON ai.id = t.item_id
		WHERE t.buyer_id = $1
		ORDER BY a.auction_date DESC, a.created_at DESC
	`, buyerID)
	if err != nil {
		return nil, dserrors.HandleError(err, "Auction", buyerID, "ListAuctionsByBuyerID")
	}
	defer func() { _ = rows.Close() }()

	var auctions []model.Auction
	for rows.Next() {
		var a model.Auction
		var auctionDate time.Time
		var startTime, endTime *time.Time
		if err := rows.Scan(
			&a.ID,
			&a.VenueID,
			&auctionDate,
			&startTime,
			&endTime,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		a.Period = model.NewAuctionPeriod(auctionDate, startTime, endTime)
		auctions = append(auctions, a)
	}
	return auctions, dserrors.HandleError(rows.Err(), "Auction", buyerID, "ListAuctionsByBuyerID")
}
