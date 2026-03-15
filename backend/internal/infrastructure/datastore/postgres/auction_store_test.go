package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/stretchr/testify/assert"
)

func TestAuctionStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewAuctionStore(postgres.NewClient(db))
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	start := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	auction := &model.Auction{
		VenueID: 1,
		Status:  model.AuctionStatusScheduled,
		Period:  model.NewAuctionPeriod(date, &start, &end),
	}

	mock.ExpectQuery("INSERT INTO auctions").
		WithArgs(auction.VenueID, auction.Period.AuctionDate, auction.Period.StartAt, auction.Period.EndAt, auction.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id", "venue_id", "auction_date", "start_time", "end_time", "status", "created_at", "updated_at"}).
			AddRow(1, 1, date, start, end, "scheduled", time.Now(), time.Now()))

	created, err := repo.Create(context.Background(), auction)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
}

func TestAuctionStore_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewAuctionStore(postgres.NewClient(db))
	id := 1
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	start := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at FROM auctions WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "venue_id", "auction_date", "start_time", "end_time", "status", "created_at", "updated_at"}).
			AddRow(1, 1, date, start, end, "scheduled", time.Now(), time.Now()))

	got, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, got.ID)
}

func TestAuctionStore_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewAuctionStore(postgres.NewClient(db))
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	start := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	t.Run("NoFilters", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at FROM auctions ORDER BY auction_date DESC, created_at DESC").
			WillReturnRows(sqlmock.NewRows([]string{"id", "venue_id", "auction_date", "start_time", "end_time", "status", "created_at", "updated_at"}).
				AddRow(1, 1, date, start, end, "scheduled", time.Now(), time.Now()))

		list, err := repo.List(context.Background(), nil)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("WithFilters", func(t *testing.T) {
		venueID := 1
		filters := &repository.AuctionFilters{VenueID: &venueID}
		mock.ExpectQuery("SELECT .* FROM auctions WHERE venue_id = \\$1 ORDER BY auction_date DESC, created_at DESC").
			WithArgs(venueID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "venue_id", "auction_date", "start_time", "end_time", "status", "created_at", "updated_at"}).
				AddRow(1, 1, date, start, end, "scheduled", time.Now(), time.Now()))

		list, err := repo.List(context.Background(), filters)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestAuctionStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewAuctionStore(postgres.NewClient(db))
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	start := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	auction := &model.Auction{
		ID:      1,
		VenueID: 1,
		Status:  model.AuctionStatusCompleted,
		Period:  model.NewAuctionPeriod(date, &start, &end),
	}

	mock.ExpectExec("UPDATE auctions SET").
		WithArgs(auction.VenueID, auction.Period.AuctionDate, auction.Period.StartAt, auction.Period.EndAt, auction.Status, auction.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(context.Background(), auction)
	assert.NoError(t, err)
}

func TestAuctionStore_UpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewAuctionStore(postgres.NewClient(db))
	id := 1
	status := model.AuctionStatusCompleted

	mock.ExpectExec("UPDATE auctions SET status = \\$1").
		WithArgs(status, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateStatus(context.Background(), id, status)
	assert.NoError(t, err)
}

func TestAuctionStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() { _ = db.Close() }()

	repo := postgres.NewAuctionStore(postgres.NewClient(db))
	id := 1

	mock.ExpectExec("DELETE FROM auctions WHERE id = \\$1").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), id)
	assert.NoError(t, err)
}
