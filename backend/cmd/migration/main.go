package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/migration"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSet("migration", flag.ExitOnError)
	cmd := fs.String("cmd", "up", "migration command (up)")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	cfg, err := config.LoadMigrationConfig()
	if err != nil {
		return fmt.Errorf("failed to load migration config: %w", err)
	}

	db, err := migration.Connect(context.Background(), cfg.DBConnectionURL())
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	switch *cmd {
	case "up":
		return migration.Up(db)
	default:
		return fmt.Errorf("unsupported command: %s", *cmd)
	}
}
