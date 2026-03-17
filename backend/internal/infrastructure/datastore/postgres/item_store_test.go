package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemStore_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewItemStore(postgres.NewClient(db))
	id := 1

	mock.ExpectQuery("(?s)SELECT .* FROM auction_items ai .*").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "auction_id", "fisherman_id", "fish_type", "quantity", "unit", "status", "created_at", "sort_order",
			"highest_bid", "highest_bidder_id", "highest_bidder_name",
		}).AddRow(id, 1, 1, "DB Tuna", 10, "kg", "Sold", time.Now(), 1, nil, nil, nil))

	item, err := repo.FindByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, "DB Tuna", item.FishType)
}

func TestItemStore_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewItemStore(postgres.NewClient(db))
	id := 1
	status := model.ItemStatusSold

	mock.ExpectExec("UPDATE auction_items SET status = \\$1 WHERE id = \\$2").
		WithArgs(status, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateStatus(context.Background(), id, status)
	assert.NoError(t, err)
}
