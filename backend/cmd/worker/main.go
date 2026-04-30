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
	"github.com/seka/fish-auction/backend/internal/worker"
	"github.com/seka/fish-auction/backend/internal/worker/handler"
)

const (
	shouldMigrate = false
	isWorker      = true
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
	repoReg, err := registry.NewRepositoryRegistry(cfg, cfg, config.NoCacheConfig, config.NoSessionConfig, shouldMigrate)
	if err != nil {
		return err
	}
	defer func() { _ = repoReg.Cleanup() }()

	// Initialize Service Registry
	serviceReg, err := registry.NewServiceRegistry(cfg, cfg, cfg, isWorker)
	if err != nil {
		return fmt.Errorf("failed to initialize service registry: %w", err)
	}

	// Create Worker
	pushRepo := repoReg.NewPushRepository()
	pushSvc := serviceReg.NewPushNotificationService()
	pushHandler := handler.NewPushNotificationHandler(pushRepo, pushSvc)

	buyerEmailSvc := serviceReg.NewBuyerEmailService()
	adminEmailSvc := serviceReg.NewAdminEmailService()
	emailHandler := handler.NewEmailHandler(buyerEmailSvc, adminEmailSvc)

	queue := serviceReg.NewJobQueue()
	w := worker.NewWorker(
		queue,
		worker.HandlerFunc(emailHandler.Handle),
		worker.HandlerFunc(pushHandler.Handle),
	)

	// Start Worker with modern signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Worker initialized. Starting...")
	return w.Start(ctx)
}
