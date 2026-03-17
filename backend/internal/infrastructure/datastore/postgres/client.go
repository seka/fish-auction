package postgres

import (
	"context"
	"database/sql"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

// Client implements datastore.Database using sql.DB.
type Client struct {
	db *sql.DB
}

var _ datastore.Database = (*Client)(nil)

// NewClient creates a new Client instance.
func NewClient(db *sql.DB) *Client {
	return &Client{db: db}
}

// Query executes a query that returns multiple rows.
func (d *Client) Query(ctx context.Context, query string, args ...any) (datastore.Rows, error) {
	var rows *sql.Rows
	var err error

	if tx, ok := GetTx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = d.db.QueryContext(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}
	return &rowsWrapper{rows: rows}, nil
}

// QueryRow executes a query that is expected to return at most one row.
func (d *Client) QueryRow(ctx context.Context, query string, args ...any) datastore.Row {
	if tx, ok := GetTx(ctx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

// Execute executes a query without returning any rows.
func (d *Client) Execute(ctx context.Context, query string, args ...any) (int64, error) {
	var res sql.Result
	var err error

	if tx, ok := GetTx(ctx); ok {
		res, err = tx.ExecContext(ctx, query, args...)
	} else {
		res, err = d.db.ExecContext(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// TransactionManager returns a new transaction manager.
func (d *Client) TransactionManager() repository.TransactionManager {
	return NewTransactionManager(d.db)
}

// Close closes the database connection.
func (d *Client) Close() error {
	return d.db.Close()
}

type rowsWrapper struct {
	rows *sql.Rows
}

func (w *rowsWrapper) Next() bool {
	return w.rows.Next()
}

func (w *rowsWrapper) Scan(dest ...any) error {
	return w.rows.Scan(dest...)
}

func (w *rowsWrapper) Close() error {
	return w.rows.Close()
}

func (w *rowsWrapper) Err() error {
	return w.rows.Err()
}
