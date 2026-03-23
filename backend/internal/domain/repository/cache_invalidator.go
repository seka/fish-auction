package repository

import "context"

// CacheInvalidator provides CacheInvalidator related functionality.
type CacheInvalidator interface {
	InvalidateCache(ctx context.Context, id int) error
}
