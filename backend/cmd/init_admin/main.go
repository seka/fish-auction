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
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		log.Println("Usage: DB_HOST=... go run cmd/init_admin/main.go")
		os.Exit(1)
	}

	dbURL := cfg.DBConnectionURL()
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	repo := postgres.NewAdminRepository(db)
	uc := admin.NewCreateAdminUseCase(repo)
	ctx := context.Background()

	count, err := uc.Count(ctx)
	if err != nil {
		log.Fatalf("Failed to count admins: %v", err)
	}

	if count > 0 {
		fmt.Printf("Admin user(s) found (%d). Skipping initialization.\n", count)
		return
	}

	fmt.Println("No admin users found. Creating default admin...")

	email := "admin@example.com"
	password := "admin-password"

	if err := uc.Execute(ctx, email, password); err != nil {
		log.Fatalf("Failed to create admin: %v", err)
	}

	fmt.Printf("Successfully created admin user: %s\n", email)
}
