package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/registry"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Worker Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// Load Config
	cfg, err := config.LoadWorkerConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize Repository Registry
	repoReg, err := registry.NewRepositoryRegistry(cfg, cfg, config.NoCacheConfig, config.NoSessionConfig)
	if err != nil {
		return err
	}
	defer func() { _ = repoReg.Cleanup() }()

	// Initialize Service Registry
	serviceReg, err := registry.NewServiceRegistry(config.NoEmailConfig, cfg, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize service registry: %w", err)
	}

	// Initialize Worker Registry
	workerReg, err := registry.NewWorkerRegistry(cfg, repoReg, serviceReg)
	if err != nil {
		return fmt.Errorf("failed to initialize worker registry: %w", err)
	}

	// Create Worker
	w, err := workerReg.NewWorker()
	if err != nil {
		return fmt.Errorf("failed to create worker: %w", err)
	}

	// Start Worker with modern signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Worker initialized. Starting...")
	return w.Start(ctx)
}
