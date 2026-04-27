package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq" // Force lib/pq driver
	"github.com/seka/fish-auction/backend/internal/domain/model"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/stretchr/testify/assert"
)

// TestVenueStore_Delete_Conflict_Integration tests actual DB behavior
// usage: go test -v -tags=integration ./internal/infrastructure/datastore/postgres/venue_delete_integration_test.go
func TestVenueStore_Delete_Conflict_Integration(t *testing.T) {
	requireIntegrationTests(t)

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

	repo := postgres.NewVenueStore(postgres.NewClient(db))

	// 2. Setup Data
	ctx := context.Background()

	// Force cleanup
	_, _ = db.ExecContext(ctx, "TRUNCATE venues, auctions, fishermen, buyers, auction_items, transactions CASCADE")

	// 2a. Create Venue
	venue := &model.Venue{Name: "Integration Test Venue", Location: "Tmp", Description: "Desc"}
	createdVenue, err := repo.Create(ctx, venue)
	assert.NoError(t, err)

	// 2b. Create Auction
	var auctionID int
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
	startTime := today.Add(9 * time.Hour)
	endTime := today.Add(21 * time.Hour)

	err = db.QueryRowContext(ctx, `
		INSERT INTO auctions (venue_id, status, start_at, end_at)
		VALUES ($1, 'scheduled', $2, $3) RETURNING id
	`, createdVenue.ID, startTime, endTime).Scan(&auctionID)
	assert.NoError(t, err)

	// 2c. Create Fisherman
	var fishermanID int
	err = db.QueryRowContext(ctx, "INSERT INTO public.fishermen (name) VALUES ('Test Fisherman') RETURNING id").Scan(&fishermanID)
	assert.NoError(t, err)

	// 2d. Create Item (linked to Auction)
	var itemID int
	t.Logf("DEBUG: venue_delete: fishermanID=%d, auctionID=%d", fishermanID, auctionID)
	err = db.QueryRowContext(ctx, `
		INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, status)
		VALUES ($1, $2, 'Maguro', 10, 'kg', 'Pending') RETURNING id
	`, fishermanID, auctionID).Scan(&itemID)
	assert.NoError(t, err)

	// 2e. Create Buyer
	var buyerID int
	err = db.QueryRowContext(ctx, "INSERT INTO public.buyers (name, organization, contact_info) VALUES ('Test Buyer', '', '') RETURNING id").Scan(&buyerID)
	assert.NoError(t, err)

	// 2f. Create Transaction (Linked to Item) -> This should BLOCK delete
	_, err = db.ExecContext(ctx, `
		INSERT INTO transactions (item_id, buyer_id, price)
		VALUES ($1, $2, 1000)
	`, itemID, buyerID)
	assert.NoError(t, err)

	// 3. Attempt Delete Venue (Logical)
	err = repo.Delete(ctx, createdVenue.ID)
	assert.NoError(t, err)

	// 4. Verify Logical Delete
	var deletedAt sql.NullTime
	err = db.QueryRowContext(ctx, "SELECT deleted_at FROM venues WHERE id = $1", createdVenue.ID).Scan(&deletedAt)
	assert.NoError(t, err)
	assert.True(t, deletedAt.Valid, "Venue should have deleted_at set")

	// 5. Cleanup
	_, _ = db.ExecContext(context.Background(), "DELETE FROM transactions WHERE item_id = $1", itemID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM auction_items WHERE id = $1", itemID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM auctions WHERE id = $1", auctionID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM fishermen WHERE id = $1", fishermanID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM buyers WHERE id = $1", buyerID)
	_, _ = db.ExecContext(context.Background(), "DELETE FROM venues WHERE id = $1", createdVenue.ID)
}
