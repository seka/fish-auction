package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
)

//go:embed seed.sql
var seedSQL string

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// Check APP_ENV explicitly
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		return fmt.Errorf("APP_ENV environment variable is required")
	}

	// Load Config
	cfg, err := config.LoadAppServerConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Safety check: Only run in development
	if cfg.AppEnv != "development" && cfg.AppEnv != "test" {
		return fmt.Errorf("seed command is only allowed in 'development' or 'test' environments. Current environment: %s", cfg.AppEnv)
	}

	ctx := context.Background()

	// Connect to DB
	connStr := cfg.DBConnectionURL()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to database. Environment:", cfg.AppEnv)

	// Clear Database
	fmt.Println("Clearing database...")
	tables := []string{
		"transactions",
		"auction_items",
		"auctions",
		"venues",
		"authentications",
		"buyers",
		"fishermen",
	}

	for _, table := range tables {
		_, err := db.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// Ignore error if table doesn't exist, but log it
			log.Printf("Warning: failed to truncate table %s: %v", table, err)
		}
	}
	fmt.Println("Database cleared.")

	// Run Seed Migration
	fmt.Println("Seeding database...")
	_, err = db.ExecContext(ctx, seedSQL)
	if err != nil {
		return fmt.Errorf("failed to execute seed SQL: %w", err)
	}

	fmt.Println("Database seeded successfully!")
	fmt.Println("Default buyer password: 'Password123'")
	return nil
}
