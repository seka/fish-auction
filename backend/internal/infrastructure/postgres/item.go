package postgres

import (
	"context"
	"database/sql"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) repository.ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	var e entity.AuctionItem
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO auction_items (fisherman_id, fish_type, quantity, unit, status) VALUES ($1, $2, $3, $4, 'Pending') RETURNING id, fisherman_id, fish_type, quantity, unit, status, created_at",
		item.FishermanID, item.FishType, item.Quantity, item.Unit,
	).Scan(&e.ID, &e.FishermanID, &e.FishType, &e.Quantity, &e.Unit, &e.Status, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return e.ToModel(), nil
}

func (r *ItemRepository) List(ctx context.Context, status string) ([]model.AuctionItem, error) {
	query := "SELECT id, fisherman_id, fish_type, quantity, unit, status, created_at FROM auction_items"
	var args []interface{}
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.AuctionItem
	for rows.Next() {
		var e entity.AuctionItem
		if err := rows.Scan(&e.ID, &e.FishermanID, &e.FishType, &e.Quantity, &e.Unit, &e.Status, &e.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, *e.ToModel())
	}
	return items, nil
}

func (r *ItemRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE auction_items SET status = $1 WHERE id = $2", status, id)
	return err
}
