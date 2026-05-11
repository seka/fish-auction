package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/seka/fish-auction/backend/config"
	apperrors "github.com/seka/fish-auction/backend/internal/domain/errors"
	"github.com/seka/fish-auction/backend/internal/infrastructure/datastore/postgres"
	"github.com/seka/fish-auction/backend/internal/logger"
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

	logger.Init(config.GetLogLevel())

	if err := run(); err != nil {
		slog.Error("init_admin fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	if email == "" || password == "" {
		return fmt.Errorf("--email and --password are required. Usage: go run cmd/init_admin/main.go --email <email> --password <password>")
	}

	cfg := config.NewAppServerConfig()
	if err := config.ValidateAppServerConfig(cfg); err != nil {
		slog.Info("usage hint", "msg", "POSTGRES_HOST=... go run cmd/init_admin/main.go --email <email> --password <password>")
		return fmt.Errorf("invalid config: %w", err)
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
