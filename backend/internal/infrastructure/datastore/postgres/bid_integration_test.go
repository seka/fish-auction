package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	cache "github.com/seka/fish-auction/backend/internal/infrastructure/datastore/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestDBConnStr() string {
	return "postgres://postgres:postgres@localhost:5432/fish_auction?sslmode=disable"
}

// TestItemStore_FindByID_IncludesHighestBid tests if FindByID returns the highest bid info.
// This is critical for CreateBidUseCase validation.
func TestItemStore_FindByID_IncludesHighestBid(t *testing.T) {
	// 1. Connect to DB
	connStr := getTestDBConnStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skip("Skipping integration test: DB connection failed: ", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.PingContext(context.Background()); err != nil {
		t.Skip("Skipping integration test: DB ping failed: ", err)
	}

	// 2. Setup Redis (Mock or Real)
	// Using real redis from docker
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	itemCache := cache.NewItemStore(cache.NewClient(redisClient), 1*time.Minute)

	repo := postgres.NewItemStore(postgres.NewClient(db))

	// Setup Data
	ctx := context.Background()

	// Force cleanup
	_, _ = db.ExecContext(ctx, "TRUNCATE public.venues, public.auctions, public.fishermen, public.buyers, public.auction_items, public.transactions CASCADE")

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	require.NoError(t, err)
	defer func() { _ = tx.Rollback() }()

	// Venue
	var venueID int
	err = tx.QueryRowContext(ctx, "INSERT INTO public.venues (name) VALUES ('Bid Test Venue') RETURNING id").Scan(&venueID)
	require.NoError(t, err)

	// Auction
	var auctionID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO public.auctions (venue_id, status, start_time, end_time, auction_date)
		VALUES ($1, 'scheduled', $2, $3, $4) RETURNING id
	`, venueID, time.Now(), time.Now().Add(1*time.Hour), time.Now()).Scan(&auctionID)
	require.NoError(t, err)

	// User (Fisherman)
	var fishermanID int
	err = tx.QueryRowContext(ctx, "INSERT INTO public.users (name, role, created_at) VALUES ('Bid Test Fisherman', 'FISHERMAN', CURRENT_TIMESTAMP) RETURNING id").Scan(&fishermanID)
	require.NoError(t, err)

	// Fisherman Profile (optional but good practice if table exists)
	_, _ = tx.ExecContext(ctx, "INSERT INTO public.fishermen (id, name) VALUES ($1, 'Bid Test Fisherman')", fishermanID)

	// Item
	var itemID int
	t.Logf("DEBUG: bid_integration: fishermanID=%d, auctionID=%d", fishermanID, auctionID)
	err = tx.QueryRowContext(ctx, `
		INSERT INTO public.auction_items (fisherman_id, auction_id, fish_type, quantity, unit, status, sort_order)
		VALUES ($1, $2, 'Katsuo', 10, 'kg', 'Pending', 1) RETURNING id
	`, fishermanID, auctionID).Scan(&itemID)
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	// Clear cache for this new item just in case
	_ = itemCache.Delete(ctx, itemID)

	// User (Buyer)
	var buyerID int
	err = db.QueryRowContext(ctx, "INSERT INTO public.users (name, role, created_at) VALUES ('Bid Test Buyer', 'BUYER', CURRENT_TIMESTAMP) RETURNING id").Scan(&buyerID)
	require.NoError(t, err)

	// Buyer profile
	_, err = db.ExecContext(ctx, "INSERT INTO public.buyers (id, name) VALUES ($1, 'Bid Test Buyer')", buyerID)
	assert.NoError(t, err)

	// Transaction (The First Bid: 1000)
	_, err = db.ExecContext(ctx, `
		INSERT INTO transactions (item_id, buyer_id, price)
		VALUES ($1, $2, 1000)
	`, itemID, buyerID)
	assert.NoError(t, err)

	// 3. Test FindByID
	item, err := repo.FindByID(ctx, itemID)
	require.NoError(t, err)
	require.NotNil(t, item)

	// 4. Assert HighestBid
	// For this specific test, we expect a highest bid of 1000.
	// The provided switch statement uses `tc.expectedPrice`, which implies a test case struct.
	// To make this change faithfully without introducing `tc` or breaking the test,
	// I will adapt the switch to directly check for 1000.
	switch {
	case item.HighestBid == nil:
		t.Errorf("expected highest bid %d, got nil", 1000)
	case item.HighestBid.Amount() != 1000:
		t.Errorf("expected highest bid %d, got %d", 1000, item.HighestBid.Amount())
	default:
		t.Logf("SUCCESS: HighestBid is %d", item.HighestBid.Amount())
	}

	// Cleanup
	_, _ = db.ExecContext(context.Background(), "DELETE FROM transactions WHERE item_id = $1", itemID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM auction_items WHERE id = $1", itemID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM fishermen WHERE id = $1", fishermanID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM auctions WHERE id = $1", auctionID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM venues WHERE id = $1", venueID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM buyers WHERE id = $1", buyerID)
}

func TestItemStore_FindByID_NoBids(t *testing.T) {
	// 1. Connect
	connStr := getTestDBConnStr()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skip("Skipping integration test: DB connection failed: ", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.PingContext(context.Background()); err != nil {
		t.Skip("Skipping integration test: DB ping failed: ", err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	itemCache := cache.NewItemStore(cache.NewClient(redisClient), 1*time.Minute)
	repo := postgres.NewItemStore(postgres.NewClient(db))
	ctx := context.Background()

	// Force cleanup
	_, _ = db.ExecContext(ctx, "TRUNCATE venues, auctions, fishermen, buyers, auction_items, transactions CASCADE")
	var venueID int
	_ = db.QueryRowContext(ctx, "INSERT INTO venues (name) VALUES ('NoBid Venue') RETURNING id").Scan(&venueID)

	var auctionID int
	_ = db.QueryRowContext(ctx, "INSERT INTO auctions (venue_id, status, start_time, end_time, auction_date) VALUES ($1, 'scheduled', $2, $3, $4) RETURNING id", venueID, time.Now(), time.Now().Add(1*time.Hour), time.Now()).Scan(&auctionID)

	// User (Fisherman)
	var fishermanID int
	err = db.QueryRowContext(ctx, "INSERT INTO public.users (name, role, created_at) VALUES ('NoBid Fisherman', 'FISHERMAN', CURRENT_TIMESTAMP) RETURNING id").Scan(&fishermanID)
	require.NoError(t, err)

	// Fisherman profile (optional)
	_, _ = db.ExecContext(ctx, "INSERT INTO public.fishermen (id, name) VALUES ($1, 'NoBid Fisherman')", fishermanID)

	var itemID int
	err = db.QueryRowContext(ctx, "INSERT INTO public.auction_items (fisherman_id, auction_id, fish_type, quantity, unit, status, sort_order) VALUES ($1, $2, 'Iwashi', 50, 'kg', 'Pending', 1) RETURNING id", fishermanID, auctionID).Scan(&itemID)
	require.NoError(t, err)

	// Clear cache
	_ = itemCache.Delete(ctx, itemID)

	// 3. Test FindByID (Expectation: Success, HighestBid is nil)
	item, err := repo.FindByID(ctx, itemID)

	if err != nil {
		t.Logf("FAILURE: FindByID returned error for item with no bids: %v", err)
		t.Fail()
	} else {
		assert.NotNil(t, item)
		if item.HighestBid != nil {
			t.Logf("FAILURE: Expected nil HighestBid, got %d", item.HighestBid.Amount())
			t.Fail()
		}
	}

	// Cleanup
	_, _ = db.ExecContext(context.Background(), "DELETE FROM auction_items WHERE id = $1", itemID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM fishermen WHERE id = $1", fishermanID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM auctions WHERE id = $1", auctionID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM venues WHERE id = $1", venueID)
}
