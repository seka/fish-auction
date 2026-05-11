package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/infrastructure/queue/sqs"
	"github.com/seka/fish-auction/backend/internal/logger"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/relay"
)

func main() {
	cfg := config.NewRelayConfig()
	logger.Init(config.GetLogLevel())

	if err := run(cfg); err != nil {
		slog.Error("relay fatal", "err", err)
		os.Exit(1)
	}
}

func run(cfg *config.RelayConfig) error {
	repoReg, err := registry.NewRepositoryRegistry(
		cfg,
		config.NoRedisConfig,
		config.NoCacheConfig,
		config.NoSessionConfig,
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

	c := relay.NewOutboxCleaner(
		outboxRepo,
		7*24*time.Hour,
		1*time.Hour,
		5*time.Minute,
		1*time.Minute,
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		r.Run(ctx)
	}()
	go func() {
		defer wg.Done()
		c.Run(ctx)
	}()

	slog.Info("relay process started", "instance_id", instanceID)

	wg.Wait()
	slog.Info("relay process stopped", "instance_id", instanceID)
	return nil
}
