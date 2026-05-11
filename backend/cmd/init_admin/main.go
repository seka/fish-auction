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
	logger.Init(config.GetLogLevel())

	flag.Parse()

	if err := run(); err != nil {
		slog.Error("init_admin fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	if email == "" || password == "" {
		return fmt.Errorf("--email and --password are required. Usage: go run cmd/init_admin/main.go --email <email> --password <password>")
	}

	if err := requireDBEnv(); err != nil {
		return err
	}

	cfg := config.NewInitAdminConfig()

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

// requireDBEnv は DB 接続に必須な環境変数の不足を起動前に検出する。
// NewInitAdminConfig は空文字をデフォルト値として返すため、未設定時に DB 接続段階で
// 失敗するまで原因が分かりづらいことを避ける。
func requireDBEnv() error {
	required := []string{"POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"}
	var missing []string
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %v. Usage: POSTGRES_HOST=... POSTGRES_PORT=... go run cmd/init_admin/main.go --email <email> --password <password>", missing)
	}
	return nil
}
