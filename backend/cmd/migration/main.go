package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/logger"
	"github.com/seka/fish-auction/backend/internal/migration"
)

const usage = `Usage: migration <command>

Commands:
  up    Apply all pending migrations
`

func main() {
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	logger.Init(config.GetLogLevel())

	cfg, err := config.LoadMigrationConfig()
	if err != nil {
		slog.Error("failed to load migration config", "err", err)
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := dispatch(cfg, flag.Arg(0), flag.Args()[1:]); err != nil {
		slog.Error("migration fatal", "err", err)
		os.Exit(1)
	}
}

func dispatch(cfg *config.MigrationConfig, subcommand string, args []string) error {
	switch subcommand {
	case "up":
		return runUp(cfg, args)
	case "help":
		flag.Usage()
		return nil
	default:
		flag.Usage()
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}

func runUp(cfg *config.MigrationConfig, args []string) error {
	fs := flag.NewFlagSet("up", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	db, err := migration.Connect(context.Background(), cfg.DBConnectionURL())
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	return migration.Up(db)
}
