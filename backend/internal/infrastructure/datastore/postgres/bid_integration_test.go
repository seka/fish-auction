package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	cache "github.com/seka/fish-auction/backend/internal/infrastructure/cache/redis"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestItemRepository_FindByID_IncludesHighestBid tests if FindByID returns the highest bid info.
// This is critical for CreateBidUseCase validation.
func TestItemRepository_FindByID_IncludesHighestBid(t *testing.T) {
	// 1. Connect to DB
	connStr := "postgres://postgres:postgres@localhost:5432/fish_auction?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skip("Skipping integration test: DB connection failed: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skip("Skipping integration test: DB ping failed: ", err)
	}

	// 2. Setup Redis (Mock or Real)
	// Using real redis from docker
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	itemCache := cache.NewItemCache(redisClient, 1*time.Minute)

	repo := postgres.NewItemRepository(db, itemCache)

	// Setup Data
	ctx := context.Background()

	// Ensure cache is clear for this itemID before test
	// We don't know itemID yet, so we'll clear after creation or just rely on new ID.

	// Venue
	var venueID int
	err = db.QueryRow("INSERT INTO venues (name) VALUES ('Bid Test Venue') RETURNING id").Scan(&venueID)
	assert.NoError(t, err)

	// Auction
	var auctionID int
	err = db.QueryRow(`
		INSERT INTO auctions (venue_id, status, start_time, end_time, auction_date) 
		VALUES ($1, 'scheduled', $2, $3, $4) RETURNING id
	`, venueID, time.Now(), time.Now().Add(1*time.Hour), time.Now()).Scan(&auctionID)
	assert.NoError(t, err)

	// Fisherman
	var fishermanID int
	err = db.QueryRow("INSERT INTO fishermen (name) VALUES ('Bid Test Fisherman') RETURNING id").Scan(&fishermanID)
	assert.NoError(t, err)

	// Item
	var itemID int
	err = db.QueryRow(`
		INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, status, sort_order)
		VALUES ($1, $2, 'Katsuo', 10, 'kg', 'Pending', 1) RETURNING id
	`, fishermanID, auctionID).Scan(&itemID)
	require.NoError(t, err)

	// Clear cache for this new item just in case
	itemCache.Delete(ctx, itemID)

	// Buyer
	var buyerID int
	err = db.QueryRow("INSERT INTO buyers (name) VALUES ('Bid Test Buyer') RETURNING id").Scan(&buyerID)
	assert.NoError(t, err)

	// Transaction (The First Bid: 1000)
	_, err = db.Exec(`
		INSERT INTO transactions (item_id, buyer_id, price)
		VALUES ($1, $2, 1000)
	`, itemID, buyerID)
	assert.NoError(t, err)

	// 3. Test FindByID
	item, err := repo.FindByID(ctx, itemID)
	require.NoError(t, err)
	require.NotNil(t, item)

	// 4. Assert HighestBid
	if item.HighestBid == nil {
		t.Log("FAILURE REPRODUCED: HighestBid is nil")
		t.Fail()
	} else if *item.HighestBid != 1000 {
		t.Logf("FAILURE: HighestBid expected 1000, got %d", *item.HighestBid)
		t.Fail()
	} else {
		t.Logf("SUCCESS: HighestBid is %d", *item.HighestBid)
	}

	// Cleanup
	db.Exec("DELETE FROM transactions WHERE item_id = $1", itemID)
	db.Exec("DELETE FROM auction_items WHERE id = $1", itemID)
	db.Exec("DELETE FROM fishermen WHERE id = $1", fishermanID)
	db.Exec("DELETE FROM auctions WHERE id = $1", auctionID)
	db.Exec("DELETE FROM venues WHERE id = $1", venueID)
	db.Exec("DELETE FROM buyers WHERE id = $1", buyerID)
}

func TestItemRepository_FindByID_NoBids(t *testing.T) {
	// 1. Connect
	connStr := "postgres://postgres:postgres@localhost:5432/fish_auction?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skip("Skipping integration test: DB connection failed: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skip("Skipping integration test: DB ping failed: ", err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	itemCache := cache.NewItemCache(redisClient, 1*time.Minute)
	repo := postgres.NewItemRepository(db, itemCache)
	ctx := context.Background()

	// 2. Setup Data (No Transaction this time)
	var venueID int
	db.QueryRow("INSERT INTO venues (name) VALUES ('NoBid Venue') RETURNING id").Scan(&venueID)

	var auctionID int
	db.QueryRow("INSERT INTO auctions (venue_id, status, start_time, end_time, auction_date) VALUES ($1, 'scheduled', $2, $3, $4) RETURNING id", venueID, time.Now(), time.Now().Add(1*time.Hour), time.Now()).Scan(&auctionID)

	var fishermanID int
	db.QueryRow("INSERT INTO fishermen (name) VALUES ('NoBid Fisherman') RETURNING id").Scan(&fishermanID)

	var itemID int
	err = db.QueryRow("INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, status, sort_order) VALUES ($1, $2, 'Iwashi', 50, 'kg', 'Pending', 1) RETURNING id", fishermanID, auctionID).Scan(&itemID)
	require.NoError(t, err)

	// Clear cache
	itemCache.Delete(ctx, itemID)

	// 3. Test FindByID (Expectation: Success, HighestBid is nil)
	item, err := repo.FindByID(ctx, itemID)

	if err != nil {
		t.Logf("FAILURE: FindByID returned error for item with no bids: %v", err)
		t.Fail()
	} else {
		assert.NotNil(t, item)
		if item.HighestBid != nil {
			t.Logf("FAILURE: Expected nil HighestBid, got %d", *item.HighestBid)
			t.Fail()
		} else {
			t.Log("SUCCESS: HighestBid is nil as expected")
		}
	}

	// Cleanup
	db.Exec("DELETE FROM auction_items WHERE id = $1", itemID)
	db.Exec("DELETE FROM fishermen WHERE id = $1", fishermanID)
	db.Exec("DELETE FROM auctions WHERE id = $1", auctionID)
	db.Exec("DELETE FROM venues WHERE id = $1", venueID)
}
