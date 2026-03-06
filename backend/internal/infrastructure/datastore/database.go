package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// Database defines the interface for database operations
type Database interface {
	// Query executes a query and returns abstract Rows interface
	Query(ctx context.Context, query string, args ...any) (Rows, error)

	// QueryRow executes a query that is expected to return at most one row
	QueryRow(ctx context.Context, query string, args ...any) Row

	// Execute executes INSERT/UPDATE/DELETE operations and returns affected rows count
	Execute(ctx context.Context, query string, args ...any) (int64, error)

	// TransactionManager returns the transaction manager for this database
	TransactionManager() repository.TransactionManager

	// Close closes the database connection
	Close() error
}

// Rows defines the interface for iterating over query results
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// Row defines the interface for a single row result
type Row interface {
	Scan(dest ...any) error
}
