package datastore

import (
	"context"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

type Database interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	Execute(ctx context.Context, query string, args ...any) (int64, error)
	TransactionManager() repository.TransactionManager
	Close() error
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

type Row interface {
	Scan(dest ...any) error
}
