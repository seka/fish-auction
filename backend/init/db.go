package initializer

import (
	"database/sql"
	"fmt"

	"log"
	"time"

	"github.com/seka/fish-auction/backend/migrations"
)

func ConnectDB(connStr string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	// Retry connection
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			return db, nil
		}
		log.Printf("Failed to connect to DB: %v. Retrying in 2s...", err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %w", err)
}

func InitDB(db *sql.DB) error {
	// Run Migrations (Simple)
	// In a real app, use a migration tool. Here we just ensure tables exist.
	migrationSQL, err := migrations.FS.ReadFile("001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
