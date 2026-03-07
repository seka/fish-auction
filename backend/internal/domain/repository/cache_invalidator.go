package repository

import "context"

type CacheInvalidator interface {
	InvalidateCache(ctx context.Context, id int) error
}
