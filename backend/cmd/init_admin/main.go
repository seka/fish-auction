package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		log.Println("Usage: DB_HOST=... go run cmd/init_admin/main.go")
		return err
	}

	ctx := context.Background()
	dbURL := cfg.DBConnectionURL()
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %w", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	repo := postgres.NewAdminStore(postgres.NewClient(db))
	uc := admin.NewCreateAdminUseCase(repo)

	count, err := uc.Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count admins: %w", err)
	}

	if count > 0 {
		fmt.Printf("Admin user(s) found (%d). Skipping initialization.\n", count)
		return nil
	}

	email := "admin@example.com"
	password := "Admin-Password123"

	if _, err = uc.Execute(ctx, email, password); err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	fmt.Printf("Successfully created admin user: %s\n", email)
	return nil
}
