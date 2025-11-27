package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server"
	"github.com/seka/fish-auction/backend/internal/server/handler"
)

func main() {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize Repository Registry (handles DB connection, Redis connection, and migration)
	connStr := cfg.DBConnectionURL()
	repoReg, db, err := registry.NewRepositoryRegistry(connStr, cfg.RedisAddr, cfg.CacheTTL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize UseCase Registry
	useCaseReg := registry.NewUseCaseRegistry(repoReg)

	// Initialize Handlers
	handlers := buildHandlers(useCaseReg)

	// Initialize Server
	srv := server.NewServer(
		handlers.health,
		handlers.fisherman,
		handlers.buyer,
		handlers.item,
		handlers.bid,
		handlers.invoice,
		handlers.auth,
	)

	// Start Server
	if err := srv.Start(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type handlers struct {
	health    *handler.HealthHandler
	fisherman *handler.FishermanHandler
	buyer     *handler.BuyerHandler
	item      *handler.ItemHandler
	bid       *handler.BidHandler
	invoice   *handler.InvoiceHandler
	auth      *handler.AuthHandler
}

func buildHandlers(reg registry.UseCase) *handlers {
	return &handlers{
		health:    handler.NewHealthHandler(),
		fisherman: handler.NewFishermanHandler(reg),
		buyer:     handler.NewBuyerHandler(reg),
		item:      handler.NewItemHandler(reg),
		bid:       handler.NewBidHandler(reg),
		invoice:   handler.NewInvoiceHandler(reg),
		auth:      handler.NewAuthHandler(reg),
	}
}
