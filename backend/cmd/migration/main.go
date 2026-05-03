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

const usage = `Usage: migration <command>

Commands:
  up    Apply all pending migrations
`

func main() {
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := dispatch(flag.Arg(0), flag.Args()[1:]); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func dispatch(subcommand string, args []string) error {
	switch subcommand {
	case "up":
		return runUp(args)
	default:
		flag.Usage()
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}

func runUp(args []string) error {
	fs := flag.NewFlagSet("up", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
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

	return migration.Up(db)
}
