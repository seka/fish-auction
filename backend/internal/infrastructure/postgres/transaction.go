package postgres

import (
	"context"
	"database/sql"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, t *model.Transaction) (*model.Transaction, error) {
	var e entity.Transaction
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO transactions (item_id, buyer_id, price) VALUES ($1, $2, $3) RETURNING id, item_id, buyer_id, price, created_at",
		t.ItemID, t.BuyerID, t.Price,
	).Scan(&e.ID, &e.ItemID, &e.BuyerID, &e.Price, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *TransactionRepository) ListInvoices(ctx context.Context) ([]model.InvoiceItem, error) {
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
