package postgres

import (
	"context"
	"database/sql"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type bidRepository struct {
	db *sql.DB
}

func NewBidRepository(db *sql.DB) repository.BidRepository {
	return &bidRepository{db: db}
}

// getDB returns the transaction if one exists in context, otherwise returns the default DB
func (r *bidRepository) getDB(ctx context.Context) dbExecutor {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return r.db
}

func (r *bidRepository) Create(ctx context.Context, bid *model.Bid) (*model.Bid, error) {
	db := r.getDB(ctx)

	e := entity.Bid{
		ItemID:  bid.ItemID,
		BuyerID: bid.BuyerID,
		Price:   bid.Price,
	}
	if err := e.Validate(); err != nil {
		return nil, err
	}

	err := db.QueryRowContext(ctx,
		"INSERT INTO transactions (item_id, buyer_id, price) VALUES ($1, $2, $3) RETURNING id, item_id, buyer_id, price, created_at",
		bid.ItemID, bid.BuyerID, bid.Price,
	).Scan(&e.ID, &e.ItemID, &e.BuyerID, &e.Price, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *bidRepository) ListInvoices(ctx context.Context) ([]model.InvoiceItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT b.id, b.name, SUM(t.price) as total_price
		FROM transactions t
		JOIN buyers b ON t.buyer_id = b.id
		GROUP BY b.id, b.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	return invoices, nil
}

func (r *bidRepository) ListPurchasesByBuyerID(ctx context.Context, buyerID int) ([]model.Purchase, error) {
	rows, err := r.db.QueryContext(ctx, `
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
		return nil, err
	}
	defer rows.Close()

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
	return purchases, rows.Err()
}

func (r *bidRepository) ListAuctionsByBuyerID(ctx context.Context, buyerID int) ([]model.Auction, error) {
	rows, err := r.db.QueryContext(ctx, `
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
		return nil, err
	}
	defer rows.Close()

	var auctions []model.Auction
	for rows.Next() {
		var a model.Auction
		if err := rows.Scan(
			&a.ID,
			&a.VenueID,
			&a.AuctionDate,
			&a.StartTime,
			&a.EndTime,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		auctions = append(auctions, a)
	}
	return auctions, rows.Err()
}
