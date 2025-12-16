package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
)

type mockItemCache struct{}

func (m *mockItemCache) Get(ctx context.Context, id int) (*model.AuctionItem, error) {
	return nil, errors.New("cache miss")
}
func (m *mockItemCache) Set(ctx context.Context, id int, item *model.AuctionItem) error { return nil }
func (m *mockItemCache) Delete(ctx context.Context, id int) error                       { return nil }

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
}
