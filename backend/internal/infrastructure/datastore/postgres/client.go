package postgres

import (
	"context"
	"database/sql"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore"
)

type client struct {
	db *sql.DB
}

// NewClient creates a new Database implementation
func NewClient(db *sql.DB) datastore.Database {
	return &client{db: db}
}

func (d *client) Query(ctx context.Context, query string, args ...any) (datastore.Rows, error) {
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

func (d *client) QueryRow(ctx context.Context, query string, args ...any) datastore.Row {
	if tx, ok := GetTx(ctx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *client) Execute(ctx context.Context, query string, args ...any) (int64, error) {
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

func (d *client) TransactionManager() repository.TransactionManager {
	return NewTransactionManager(d.db)
}

func (d *client) Close() error {
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
