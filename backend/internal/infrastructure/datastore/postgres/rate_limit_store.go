package postgres

import (
	"context"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

var _ repository.RateLimitRepository = (*RateLimitStore)(nil)

// RateLimitStore implements repository.RateLimitRepository using PostgreSQL.
type RateLimitStore struct {
	db datastore.Database
}

// NewRateLimitStore creates a new RateLimitStore.
func NewRateLimitStore(db datastore.Database) *RateLimitStore {
	return &RateLimitStore{db: db}
}

// Increment atomically increments the counter for key within the given window.
// When the window changes, the count resets to 1.
func (s *RateLimitStore) Increment(ctx context.Context, key string, windowStart time.Time) (int64, error) {
	const q = `
		INSERT INTO rate_limit_counters (key, count, window_start)
		VALUES ($1, 1, $2)
		ON CONFLICT (key) DO UPDATE
		  SET
		    count = CASE
		      WHEN rate_limit_counters.window_start = EXCLUDED.window_start
		      THEN rate_limit_counters.count + 1
		      ELSE 1
		    END,
		    window_start = EXCLUDED.window_start
		RETURNING count`

	var count int64
	if err := s.db.QueryRow(ctx, q, key, windowStart).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
