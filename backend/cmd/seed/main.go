package main

import (
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
	// Check APP_ENV explicitly
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		log.Fatal("APP_ENV environment variable is required")
	}

	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Safety check: Only run in development
	if cfg.AppEnv != "development" && cfg.AppEnv != "test" {
		log.Fatalf("Seed command is only allowed in 'development' or 'test' environments. Current environment: %s", cfg.AppEnv)
	}

	// Connect to DB
	connStr := cfg.DBConnectionURL()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
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
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// Ignore error if table doesn't exist, but log it
			log.Printf("Warning: failed to truncate table %s: %v", table, err)
		}
	}
	fmt.Println("Database cleared.")

	// Run Seed Migration
	fmt.Println("Seeding database...")
	_, err = db.Exec(seedSQL)
	if err != nil {
		log.Fatalf("Failed to execute seed SQL: %v", err)
	}

	fmt.Println("Database seeded successfully!")
}
