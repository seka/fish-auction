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
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize Repository Registry
	repoReg, err := registry.NewRepositoryRegistry(cfg)
	if err != nil {
		return err
	}
	defer func() { _ = repoReg.Cleanup() }()

	// Initialize Service Registry
	serviceReg := registry.NewServiceRegistry(cfg)

	// Initialize Worker Registry
	workerReg := registry.NewWorkerRegistry(cfg, repoReg, serviceReg)

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
