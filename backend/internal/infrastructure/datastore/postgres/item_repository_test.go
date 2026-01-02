package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockItemCache struct {
	getFunc    func(ctx context.Context, id int) (*model.AuctionItem, error)
	setFunc    func(ctx context.Context, id int, item *model.AuctionItem) error
	deleteFunc func(ctx context.Context, id int) error
}

func (m *mockItemCache) Get(ctx context.Context, id int) (*model.AuctionItem, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return nil, errors.New("cache miss")
}
func (m *mockItemCache) Set(ctx context.Context, id int, item *model.AuctionItem) error {
	if m.setFunc != nil {
		return m.setFunc(ctx, id, item)
	}
	return nil
}
func (m *mockItemCache) Delete(ctx context.Context, id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func TestItemRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewItemRepository(db, &mockItemCache{})
	item := &model.AuctionItem{
		AuctionID:   1,
		FishermanID: 1,
		FishType:    "Tuna",
		Quantity:    10,
		Unit:        "kg",
	}

	mock.ExpectQuery("INSERT INTO auction_items").
		WithArgs(item.AuctionID, item.FishermanID, item.FishType, item.Quantity, item.Unit).
		WillReturnRows(sqlmock.NewRows([]string{"id", "auction_id", "fisherman_id", "fish_type", "quantity", "unit", "status", "created_at"}).
			AddRow(1, item.AuctionID, item.FishermanID, item.FishType, item.Quantity, item.Unit, "Pending", time.Now()))

	created, err := repo.Create(context.Background(), item)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
}

func TestItemRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewItemRepository(db, &mockItemCache{})

	t.Run("NoFilter", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM auction_items ORDER BY created_at DESC").
			WillReturnRows(sqlmock.NewRows([]string{"id", "auction_id", "fisherman_id", "fish_type", "quantity", "unit", "status", "created_at"}).
				AddRow(1, 1, 1, "Tuna", 10, "kg", "Pending", time.Now()))

		list, err := repo.List(context.Background(), "")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("WithStatus", func(t *testing.T) {
		status := "Pending"
		mock.ExpectQuery("SELECT .* FROM auction_items WHERE status = \\$1 ORDER BY created_at DESC").
			WithArgs(status).
			WillReturnRows(sqlmock.NewRows([]string{"id", "auction_id", "fisherman_id", "fish_type", "quantity", "unit", "status", "created_at"}).
				AddRow(2, 1, 1, "Salmon", 5, "kg", "Pending", time.Now()))

		list, err := repo.List(context.Background(), status)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestItemRepository_ListByAuction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewItemRepository(db, &mockItemCache{})
	auctionID := 10

	// Use (?s) to allow dot to match newlines
	mock.ExpectQuery(`(?s)SELECT .* FROM auction_items .* WHERE ai.auction_id = \$1.*`).
		WithArgs(auctionID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "auction_id", "fisherman_id", "fish_type", "quantity", "unit", "status", "created_at",
			"highest_bid", "highest_bidder_id", "highest_bidder_name",
		}).AddRow(1, auctionID, 1, "Tuna", 10, "kg", "Pending", time.Now(), 1000, 5, "Buyer A"))

	items, err := repo.ListByAuction(context.Background(), auctionID)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.NotNil(t, items[0].HighestBid)
	assert.Equal(t, 1000, *items[0].HighestBid)
}

func TestItemRepository_FindByID(t *testing.T) {
	t.Run("CacheHit", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		defer db.Close()

		mockCache := &mockItemCache{
			getFunc: func(ctx context.Context, id int) (*model.AuctionItem, error) {
				return &model.AuctionItem{ID: id, FishType: "Cached Tuna"}, nil
			},
		}

		repo := postgres.NewItemRepository(db, mockCache)
		item, err := repo.FindByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, "Cached Tuna", item.FishType)
	})

	t.Run("CacheMiss_DBHit", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		mockCache := &mockItemCache{
			getFunc: func(ctx context.Context, id int) (*model.AuctionItem, error) {
				return nil, errors.New("miss")
			},
			setFunc: func(ctx context.Context, id int, item *model.AuctionItem) error {
				return nil
			},
		}

		repo := postgres.NewItemRepository(db, mockCache)
		id := 1

		mock.ExpectQuery("SELECT .* FROM auction_items WHERE id = \\$1").
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "auction_id", "fisherman_id", "fish_type", "quantity", "unit", "status", "created_at"}).
				AddRow(id, 1, 1, "DB Tuna", 10, "kg", "Sold", time.Now()))

		item, err := repo.FindByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, "DB Tuna", item.FishType)
	})
}

func TestItemRepository_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockCache := &mockItemCache{
		deleteFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	repo := postgres.NewItemRepository(db, mockCache)
	id := 1
	status := model.ItemStatusSold

	mock.ExpectExec("UPDATE auction_items SET status = \\$1 WHERE id = \\$2").
		WithArgs(status, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateStatus(context.Background(), id, status)
	assert.NoError(t, err)
}
