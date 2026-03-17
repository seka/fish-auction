package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// Database defines the interface for low-level database operations.
type Database interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	Execute(ctx context.Context, query string, args ...any) (int64, error)
	TransactionManager() repository.TransactionManager
	Close() error
}

// Rows defines the interface for iterating over database query results.
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// Row defines the interface for a single database row result.
type Row interface {
	Scan(dest ...any) error
}
