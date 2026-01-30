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

// TestVenueRepository_Delete_Conflict_Integration tests actual DB behavior
// usage: go test -v -tags=integration ./internal/infrastructure/datastore/postgres/venue_delete_integration_test.go
func TestVenueRepository_Delete_Conflict_Integration(t *testing.T) {
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

	repo := postgres.NewVenueRepository(db)

	// 2. Setup Data
	ctx := context.Background()

	// 2a. Create Venue
	venue := &model.Venue{Name: "Integration Test Venue", Location: "Tmp", Description: "Desc"}
	createdVenue, err := repo.Create(ctx, venue)
	assert.NoError(t, err)

	// 2b. Create Auction
	var auctionID int
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(jst)
	err = db.QueryRow(`
		INSERT INTO auctions (venue_id, status, start_time, end_time, auction_date) 
		VALUES ($1, 'scheduled', $2, $3, $4) RETURNING id
	`, createdVenue.ID, now, now.Add(1*time.Hour), now).Scan(&auctionID)
	assert.NoError(t, err)

	// 2c. Create Fisherman
	var fishermanID int
	err = db.QueryRow(`INSERT INTO fishermen (name) VALUES ('Test Fisherman') RETURNING id`).Scan(&fishermanID)
	assert.NoError(t, err)

	// 2d. Create Item (linked to Auction)
	var itemID int
	err = db.QueryRow(`
		INSERT INTO auction_items (fisherman_id, auction_id, fish_type, quantity, unit, status)
		VALUES ($1, $2, 'Maguro', 10, 'kg', 'Pending') RETURNING id
	`, fishermanID, auctionID).Scan(&itemID)
	assert.NoError(t, err)

	// 2e. Create Buyer
	var buyerID int
	err = db.QueryRow(`INSERT INTO buyers (name) VALUES ('Test Buyer') RETURNING id`).Scan(&buyerID)
	assert.NoError(t, err)

	// 2f. Create Transaction (Linked to Item) -> This should BLOCK delete
	_, err = db.Exec(`
		INSERT INTO transactions (item_id, buyer_id, price)
		VALUES ($1, $2, 1000)
	`, itemID, buyerID)
	assert.NoError(t, err)

	// 3. Attempt Delete Venue (Logical)
	err = repo.Delete(ctx, createdVenue.ID)
	assert.NoError(t, err)

	// 4. Verify Logical Delete
	var deletedAt sql.NullTime
	err = db.QueryRow("SELECT deleted_at FROM venues WHERE id = $1", createdVenue.ID).Scan(&deletedAt)
	assert.NoError(t, err)
	assert.True(t, deletedAt.Valid, "Venue should have deleted_at set")

	// 5. Cleanup
	db.Exec("DELETE FROM transactions WHERE item_id = $1", itemID)
	db.Exec("DELETE FROM auction_items WHERE id = $1", itemID)
	db.Exec("DELETE FROM auctions WHERE id = $1", auctionID)
	db.Exec("DELETE FROM fishermen WHERE id = $1", fishermanID)
	db.Exec("DELETE FROM buyers WHERE id = $1", buyerID)
	db.Exec("DELETE FROM venues WHERE id = $1", createdVenue.ID)
}
