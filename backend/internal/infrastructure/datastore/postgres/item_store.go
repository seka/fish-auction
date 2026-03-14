package postgres

import (
	"context"
	"database/sql"

	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
	dserrors "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres/errors"
	"github.com/seka/fish-auction/backend/internal/infrastructure/entity"
)

type itemStore struct {
	db datastore.Database
}

var _ repository.ItemRepository = (*itemStore)(nil)

func NewItemStore(db datastore.Database) *itemStore {
	return &itemStore{db: db}
}

func (r *itemStore) Create(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	e := entity.AuctionItem{
		AuctionID:   item.AuctionID,
		FishermanID: item.FishermanID,
		FishType:    item.FishType,
		Quantity:    item.Quantity,
		Unit:        item.Unit,
		Status:      model.ItemStatusPending,
	}
	if err := e.Validate(); err != nil {
		return nil, err
	}

	err := r.db.QueryRow(ctx,
		"INSERT INTO auction_items (auction_id, fisherman_id, fish_type, quantity, unit, status) VALUES ($1, $2, $3, $4, $5, 'Pending') RETURNING id, auction_id, fisherman_id, fish_type, quantity, unit, status, sort_order, created_at",
		item.AuctionID, item.FishermanID, item.FishType, item.Quantity, item.Unit,
	).Scan(&e.ID, &e.AuctionID, &e.FishermanID, &e.FishType, &e.Quantity, &e.Unit, &e.Status, &e.SortOrder, &e.CreatedAt)
	if err != nil {
		return nil, dserrors.HandleError(err, "Item", nil, "failed to create item")
	}
	return e.ToModel(), nil
}

func (r *itemStore) List(ctx context.Context, status string) ([]model.AuctionItem, error) {
	query := "SELECT id, auction_id, fisherman_id, fish_type, quantity, unit, status, sort_order, created_at, deleted_at FROM auction_items WHERE deleted_at IS NULL"
	var args []interface{}
	if status != "" {
		query += " AND status = $1"
		args = append(args, status)
	}
	query += " ORDER BY auction_id DESC, sort_order ASC, created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, dserrors.HandleError(err, "Item", nil, "failed to list items")
	}
	defer func() { _ = rows.Close() }()

	var items []model.AuctionItem
	for rows.Next() {
		var e entity.AuctionItem
		if err := rows.Scan(&e.ID, &e.AuctionID, &e.FishermanID, &e.FishType, &e.Quantity, &e.Unit, &e.Status, &e.SortOrder, &e.CreatedAt, &e.DeletedAt); err != nil {
			return nil, dserrors.HandleError(err, "Item", nil, "failed to scan item row")
		}
		items = append(items, *e.ToModel())
	}
	if err := rows.Err(); err != nil {
		return nil, dserrors.HandleError(err, "Item", nil, "failed to iterate item rows")
	}
	return items, nil
}

func (r *itemStore) ListByAuction(ctx context.Context, auctionID int) ([]model.AuctionItem, error) {
	query := `
		SELECT
			ai.id, ai.auction_id, ai.fisherman_id, ai.fish_type,
			ai.quantity, ai.unit, ai.status, ai.created_at, ai.sort_order,
			t_max.max_price as highest_bid,
			t_max.buyer_id as highest_bidder_id,
			b.name as highest_bidder_name
		FROM auction_items ai
		LEFT JOIN (
			SELECT
				t1.item_id,
				MAX(t1.price) as max_price,
				(SELECT t2.buyer_id FROM transactions t2
				 WHERE t2.item_id = t1.item_id
				 ORDER BY t2.price DESC, t2.created_at ASC
				 LIMIT 1) as buyer_id
			FROM transactions t1
			GROUP BY t1.item_id
		) t_max ON ai.id = t_max.item_id
		LEFT JOIN buyers b ON t_max.buyer_id = b.id
		WHERE ai.auction_id = $1 AND ai.deleted_at IS NULL
		ORDER BY ai.sort_order ASC, ai.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, auctionID)
	if err != nil {
		return nil, dserrors.HandleError(err, "Item", nil, "failed to list items by auction")
	}
	defer func() { _ = rows.Close() }()

	var items []model.AuctionItem
	for rows.Next() {
		var e entity.AuctionItem
		var highestBid sql.NullInt64
		var highestBidderID sql.NullInt64
		var highestBidderName sql.NullString

		if err := rows.Scan(
			&e.ID, &e.AuctionID, &e.FishermanID, &e.FishType,
			&e.Quantity, &e.Unit, &e.Status, &e.CreatedAt,
			&e.SortOrder,
			&highestBid, &highestBidderID, &highestBidderName,
		); err != nil {
			return nil, dserrors.HandleError(err, "Item", nil, "failed to scan item row")
		}

		if highestBid.Valid {
			bid := int(highestBid.Int64)
			e.HighestBid = &bid
		}
		if highestBidderID.Valid {
			bidderID := int(highestBidderID.Int64)
			e.HighestBidderID = &bidderID
		}
		if highestBidderName.Valid {
			e.HighestBidderName = &highestBidderName.String
		}

		items = append(items, *e.ToModel())
	}
	if err := rows.Err(); err != nil {
		return nil, dserrors.HandleError(err, "Item", nil, "failed to iterate item rows")
	}
	return items, nil
}

func (r *itemStore) FindByID(ctx context.Context, id int) (*model.AuctionItem, error) {
	var e entity.AuctionItem
	var highestBid sql.NullInt64
	var highestBidderID sql.NullInt64
	var highestBidderName sql.NullString

	query := `
		SELECT
			ai.id, ai.auction_id, ai.fisherman_id, ai.fish_type,
			ai.quantity, ai.unit, ai.status, ai.created_at, ai.sort_order,
			t_max.max_price as highest_bid,
			t_max.buyer_id as highest_bidder_id,
			b.name as highest_bidder_name
		FROM auction_items ai
		LEFT JOIN (
			SELECT
				t1.item_id,
				MAX(t1.price) as max_price,
				(SELECT t2.buyer_id FROM transactions t2
				 WHERE t2.item_id = t1.item_id
				 ORDER BY t2.price DESC, t2.created_at ASC
				 LIMIT 1) as buyer_id
			FROM transactions t1
			WHERE t1.item_id = $1
			GROUP BY t1.item_id
		) t_max ON ai.id = t_max.item_id
		LEFT JOIN buyers b ON t_max.buyer_id = b.id
		WHERE ai.id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.AuctionID, &e.FishermanID, &e.FishType,
		&e.Quantity, &e.Unit, &e.Status, &e.CreatedAt,
		&e.SortOrder,
		&highestBid, &highestBidderID, &highestBidderName,
	)

	if err != nil {
		return nil, dserrors.HandleError(err, "Item", id, "failed to find item by ID")
	}

	if highestBid.Valid {
		bid := int(highestBid.Int64)
		e.HighestBid = &bid
	}
	if highestBidderID.Valid {
		bidderID := int(highestBidderID.Int64)
		e.HighestBidderID = &bidderID
	}
	if highestBidderName.Valid {
		e.HighestBidderName = &highestBidderName.String
	}

	return e.ToModel(), nil
}

func (r *itemStore) FindByIDWithLock(ctx context.Context, id int) (*model.AuctionItem, error) {
	var e entity.AuctionItem
	var highestBid sql.NullInt64
	var highestBidderID sql.NullInt64
	var highestBidderName sql.NullString

	query := `
		SELECT
			ai.id, ai.auction_id, ai.fisherman_id, ai.fish_type,
			ai.quantity, ai.unit, ai.status, ai.created_at, ai.sort_order,
			t_max.max_price as highest_bid,
			t_max.buyer_id as highest_bidder_id,
			b.name as highest_bidder_name
		FROM auction_items ai
		LEFT JOIN (
			SELECT
				t1.item_id,
				MAX(t1.price) as max_price,
				(SELECT t2.buyer_id FROM transactions t2
				 WHERE t2.item_id = t1.item_id
				 ORDER BY t2.price DESC, t2.created_at ASC
				 LIMIT 1) as buyer_id
			FROM transactions t1
			WHERE t1.item_id = $1
			GROUP BY t1.item_id
		) t_max ON ai.id = t_max.item_id
		LEFT JOIN buyers b ON t_max.buyer_id = b.id
		WHERE ai.id = $1
		FOR UPDATE OF ai
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.AuctionID, &e.FishermanID, &e.FishType,
		&e.Quantity, &e.Unit, &e.Status, &e.CreatedAt,
		&e.SortOrder,
		&highestBid, &highestBidderID, &highestBidderName,
	)

	if err != nil {
		return nil, dserrors.HandleError(err, "Item", id, "failed to find item by ID with lock")
	}

	if highestBid.Valid {
		bid := int(highestBid.Int64)
		e.HighestBid = &bid
	}
	if highestBidderID.Valid {
		bidderID := int(highestBidderID.Int64)
		e.HighestBidderID = &bidderID
	}
	if highestBidderName.Valid {
		e.HighestBidderName = &highestBidderName.String
	}

	return e.ToModel(), nil
}

func (r *itemStore) UpdateStatus(ctx context.Context, id int, status model.ItemStatus) error {
	_, err := r.db.Execute(ctx, "UPDATE auction_items SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return dserrors.HandleError(err, "Item", id, "failed to update item status")
	}
	return nil
}

func (r *itemStore) Update(ctx context.Context, item *model.AuctionItem) (*model.AuctionItem, error) {
	e := entity.AuctionItem{
		ID:          item.ID,
		AuctionID:   item.AuctionID,
		FishermanID: item.FishermanID,
		FishType:    item.FishType,
		Quantity:    item.Quantity,
		Unit:        item.Unit,
		Status:      item.Status,
	}

	if err := e.Validate(); err != nil {
		return nil, err
	}

	query := `
		UPDATE auction_items
		SET auction_id = $1, fisherman_id = $2, fish_type = $3, quantity = $4, unit = $5, status = $6
		WHERE id = $7
		RETURNING id, auction_id, fisherman_id, fish_type, quantity, unit, status, sort_order, created_at
	`
	err := r.db.QueryRow(ctx, query, e.AuctionID, e.FishermanID, e.FishType, e.Quantity, e.Unit, e.Status, e.ID).
		Scan(&e.ID, &e.AuctionID, &e.FishermanID, &e.FishType, &e.Quantity, &e.Unit, &e.Status, &e.SortOrder, &e.CreatedAt)

	if err != nil {
		return nil, dserrors.HandleError(err, "Item", e.ID, "failed to update item")
	}

	return e.ToModel(), nil
}

func (r *itemStore) Delete(ctx context.Context, id int) error {
	_, err := r.db.Execute(ctx, "UPDATE auction_items SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", id)
	if err != nil {
		return dserrors.HandleError(err, "Item", id, "failed to delete item")
	}
	return nil
}

func (r *itemStore) UpdateSortOrder(ctx context.Context, id int, sortOrder int) error {
	_, err := r.db.Execute(ctx, "UPDATE auction_items SET sort_order = $1 WHERE id = $2", sortOrder, id)
	if err != nil {
		return dserrors.HandleError(err, "Item", id, "failed to update item sort order")
	}
	return nil
}

func (r *itemStore) Reorder(ctx context.Context, auctionID int, ids []int) error {
	txMgr := r.db.TransactionManager()
	if txMgr != nil {
		return txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
			return r.doReorder(txCtx, auctionID, ids)
		})
	}
	return r.doReorder(ctx, auctionID, ids)
}

func (r *itemStore) doReorder(ctx context.Context, auctionID int, ids []int) error {
	for i, id := range ids {
		newSortOrder := i + 1
		_, err := r.db.Execute(ctx, "UPDATE auction_items SET sort_order = $1 WHERE id = $2 AND auction_id = $3", newSortOrder, id, auctionID)
		if err != nil {
			return dserrors.HandleError(err, "Item", id, "failed to update item sort order during reorder")
		}
	}
	return nil
}
