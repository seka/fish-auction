package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	domainrepo "github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server"
	"github.com/seka/fish-auction/backend/internal/server/handler"
)

type handlers struct {
	health         *handler.HealthHandler
	fisherman      *handler.FishermanHandler
	buyer          *handler.BuyerHandler
	adminBuyer     *handler.AdminBuyerHandler
	publicItem     *handler.PublicItemHandler
	adminItem      *handler.AdminItemHandler
	bid            *handler.BidHandler
	invoice        *handler.InvoiceHandler
	auth           *handler.AuthHandler
	publicVenue    *handler.PublicVenueHandler
	adminVenue     *handler.AdminVenueHandler
	publicAuction  *handler.PublicAuctionHandler
	adminAuction   *handler.AdminAuctionHandler
	admin          *handler.AdminHandler
	authReset      *handler.AuthResetHandler
	adminAuthReset *handler.AdminAuthResetHandler
	push           *handler.PushHandler
}

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize Repository Registry (handles DB connection, Redis connection, and migration)
	repoReg, db, err := registry.NewRepositoryRegistry(cfg)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	// Initialize Service Registry
	serviceReg := registry.NewServiceRegistry(cfg)

	// Initialize UseCase Registry
	useCaseReg := registry.NewUseCaseRegistry(repoReg, serviceReg)

	// Initialize Handlers
	sessionRepo := repoReg.NewSessionRepository()
	h := buildHandlers(useCaseReg, sessionRepo)

	// Initialize Server
	srv := server.NewServer(
		h.health,
		h.fisherman,
		h.buyer,
		h.adminBuyer,
		h.publicItem,
		h.adminItem,
		h.bid,
		h.invoice,
		h.auth,
		h.publicVenue,
		h.adminVenue,
		h.publicAuction,
		h.adminAuction,
		h.admin,
		h.authReset,
		h.adminAuthReset,
		h.push,
		sessionRepo,
		strings.Split(cfg.AllowedOrigins, ","),
		cfg.ReadTimeoutSec,
		cfg.WriteTimeoutSec,
		cfg.IdleTimeoutSec,
	)

	// Start Server
	if err := srv.Start(cfg.ServerAddress); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func buildHandlers(reg registry.UseCase, sessionRepo domainrepo.SessionRepository) *handlers {
	return &handlers{
		health:         handler.NewHealthHandler(),
		fisherman:      handler.NewFishermanHandler(reg),
		buyer:          handler.NewBuyerHandler(reg, sessionRepo),
		adminBuyer:     handler.NewAdminBuyerHandler(reg),
		publicItem:     handler.NewPublicItemHandler(reg),
		adminItem:      handler.NewAdminItemHandler(reg),
		bid:            handler.NewBidHandler(reg),
		invoice:        handler.NewInvoiceHandler(reg),
		auth:           handler.NewAuthHandler(reg, sessionRepo),
		publicVenue:    handler.NewPublicVenueHandler(reg),
		adminVenue:     handler.NewAdminVenueHandler(reg),
		publicAuction:  handler.NewPublicAuctionHandler(reg),
		adminAuction:   handler.NewAdminAuctionHandler(reg),
		admin:          handler.NewAdminHandler(reg),
		authReset:      handler.NewAuthResetHandler(reg),
		adminAuthReset: handler.NewAdminAuthResetHandler(reg),
		push:           handler.NewPushHandler(reg),
	}
}
