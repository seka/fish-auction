package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
)

// txKey is used to store transaction in context
type txKey struct{}

type transactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sql.DB) repository.TransactionManager {
	return &transactionManager{db: db}
}

// WithTransaction executes the given function within a database transaction
func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Store transaction in context
	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// GetTx retrieves the transaction from context, if any
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}
