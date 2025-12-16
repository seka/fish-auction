package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/stretchr/testify/assert"
)

func TestAuctionRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAuctionRepository(db)
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	start := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	auction := &model.Auction{
		VenueID:     1,
		AuctionDate: date,
		StartTime:   &start,
		EndTime:     &end,
		Status:      model.AuctionStatusScheduled,
	}

	mock.ExpectQuery("INSERT INTO auctions").
		WithArgs(auction.VenueID, auction.AuctionDate, auction.StartTime, auction.EndTime, auction.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id", "venue_id", "auction_date", "start_time", "end_time", "status", "created_at", "updated_at"}).
			AddRow(1, 1, date, start, end, "scheduled", time.Now(), time.Now()))

	created, err := repo.Create(context.Background(), auction)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
}

func TestAuctionRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAuctionRepository(db)
	id := 1
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	start := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT id, venue_id, auction_date, start_time, end_time, status, created_at, updated_at FROM auctions WHERE id = \\$1").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "venue_id", "auction_date", "start_time", "end_time", "status", "created_at", "updated_at"}).
			AddRow(1, 1, date, start, end, "scheduled", time.Now(), time.Now()))

	got, err := repo.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, got.ID)
}

func TestAuctionRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAuctionRepository(db)
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
