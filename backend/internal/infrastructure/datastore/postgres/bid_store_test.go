package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/stretchr/testify/assert"
)

func TestBidStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewBidStore(postgres.NewClient(db))
	bid := &model.Bid{
		ItemID:  101,
		BuyerID: 1,
		Price:   model.NewBidPrice(1500),
	}

	mock.ExpectQuery("INSERT INTO transactions").
		WithArgs(bid.ItemID, bid.BuyerID, bid.Price.Amount()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "item_id", "buyer_id", "price", "created_at"}).
			AddRow(1, bid.ItemID, bid.BuyerID, bid.Price.Amount(), time.Now()))

	created, err := repo.Create(context.Background(), bid)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
}

func TestBidStore_ListPurchasesByBuyerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewBidStore(postgres.NewClient(db))
	buyerID := 1

	mock.ExpectQuery("SELECT t.id, t.item_id, ai.fish_type, ai.quantity, ai.unit, t.price, t.buyer_id, ai.auction_id, TO_CHAR\\(a\\.start_at AT TIME ZONE 'Asia/Tokyo', 'YYYY-MM-DD'\\), t.created_at FROM transactions t .* WHERE t.buyer_id = \\$1").
		WithArgs(buyerID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "item_id", "fish_type", "quantity", "unit", "price", "buyer_id", "auction_id", "start_at", "created_at"}).
			AddRow(1, 101, "Tuna", 1, "kg", 1500, buyerID, 1, "2023-01-01", time.Now()))

	list, err := repo.ListPurchasesByBuyerID(context.Background(), buyerID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}
