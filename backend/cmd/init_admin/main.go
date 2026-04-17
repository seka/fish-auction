package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/seka/fish-auction/backend/config"
	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/seka/fish-auction/backend/internal/usecase/admin"

	_ "github.com/lib/pq"
)

var (
	email    string
	password string
)

func init() {
	flag.StringVar(&email, "email", "", "admin email (required)")
	flag.StringVar(&password, "password", "", "admin password (required)")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if email == "" || password == "" {
		return fmt.Errorf("--email and --password are required. Usage: go run cmd/init_admin/main.go --email <email> --password <password>")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		log.Println("Usage: POSTGRES_HOST=... go run cmd/init_admin/main.go --email <email> --password <password>")
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

	if _, err = uc.Execute(ctx, email, password); err != nil {
		var conflictErr *apperrors.ConflictError
		if errors.As(err, &conflictErr) {
			fmt.Printf("Admin user with email %s already exists. Skipping.\n", email)
			return nil
		}
		return fmt.Errorf("failed to create admin: %w", err)
	}

	fmt.Printf("Successfully created admin user: %s\n", email)
	return nil
}
