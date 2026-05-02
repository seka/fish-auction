package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/sqs"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/relay"
)

const shouldMigrate = false

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadRelayConfig()
	if err != nil {
		return fmt.Errorf("failed to load relay config: %w", err)
	}

	repoReg, err := registry.NewRepositoryRegistry(
		cfg,
		config.NoRedisConfig,
		config.NoCacheConfig,
		config.NoSessionConfig,
		shouldMigrate,
	)
	if err != nil {
		return err
	}
	defer func() { _ = repoReg.Cleanup() }()

	outboxRepo := repoReg.NewOutboxRepository()

	region, queueURL, endpoint := cfg.SQSConfig()
	jobQueue, err := sqs.NewClient(context.Background(), region, queueURL, endpoint)
	if err != nil {
		return fmt.Errorf("failed to initialize SQS client: %w", err)
	}

	hostname, _ := os.Hostname()
	instanceID := fmt.Sprintf("relay-%s-%d", hostname, os.Getpid())

	r := relay.NewOutboxRelay(
		outboxRepo,
		jobQueue,
		5*time.Second,
		50,
		instanceID,
	)

	c := relay.NewOutboxCleanup(
		outboxRepo,
		7*24*time.Hour,
		1*time.Hour,
		5*time.Minute,
		1*time.Minute,
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	r.Start(ctx, &wg)
	c.Start(ctx, &wg)

	log.Printf("Relay process started (instance=%s)", instanceID)

	<-ctx.Done()
	log.Println("Relay process shutting down...")
	wg.Wait()
	log.Println("Relay process stopped")
	return nil
}
