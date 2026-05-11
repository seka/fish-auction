package migration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/migrations"
)

const (
	connectRetries  = 10
	connectInterval = 2 * time.Second
)

// Connect establishes a *sql.DB connection with retry logic, mirroring the
// behavior previously embedded in the repository registry. The caller owns
// the returned *sql.DB and must close it.
func Connect(ctx context.Context, dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for range connectRetries {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.PingContext(ctx)
		}
		if err == nil {
			return db, nil
		}
		slog.Warn("failed to connect to DB; retrying", "err", err, "interval", connectInterval.String())
		time.Sleep(connectInterval)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %w", err)
}

// Up applies all pending migrations against the supplied database.
// Returns an error if migrations fail or the database is left in a dirty state.
func Up(db *sql.DB) error {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get migration version: %w", err)
	}
	if dirty {
		return fmt.Errorf("migration is in dirty state at version %d, manual intervention required", version)
	}
	slog.Info("migration complete", "version", version, "dirty", dirty)
	return nil
}
