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
	adminHandler "github.com/seka/fish-auction/backend/internal/server/handler/admin"
	buyerHandler "github.com/seka/fish-auction/backend/internal/server/handler/buyer"
	publicHandler "github.com/seka/fish-auction/backend/internal/server/handler/public"
)

type handlers struct {
	health         *publicHandler.HealthHandler
	fisherman      *adminHandler.FishermanHandler
	buyerAuth      *publicHandler.BuyerAuthHandler
	buyer          *buyerHandler.BuyerHandler
	adminBuyer     *adminHandler.BuyerHandler
	publicItem     *publicHandler.ItemHandler
	adminItem      *adminHandler.ItemHandler
	bid            *buyerHandler.BidHandler
	invoice        *adminHandler.InvoiceHandler
	adminAuth      *publicHandler.AdminAuthHandler
	publicVenue    *publicHandler.VenueHandler
	adminVenue     *adminHandler.VenueHandler
	publicAuction  *publicHandler.AuctionHandler
	adminAuction   *adminHandler.AuctionHandler
	admin          *adminHandler.AdminHandler
	authReset      *publicHandler.AuthResetHandler
	adminAuthReset *adminHandler.AuthResetHandler
	push           *buyerHandler.PushHandler
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
	useCaseReg := registry.NewUseCaseRegistry(repoReg, serviceReg, cfg)

	// Initialize Handlers
	sessionRepo := repoReg.NewSessionRepository()
	h := buildHandlers(useCaseReg, sessionRepo)

	// Initialize Server
	srv := server.NewServer(
		h.health,
		h.fisherman,
		h.buyerAuth,
		h.buyer,
		h.adminBuyer,
		h.publicItem,
		h.adminItem,
		h.bid,
		h.invoice,
		h.adminAuth,
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
		cfg.ReadTimeout,
		cfg.WriteTimeout,
		cfg.IdleTimeout,
	)

	// Start Server
	if err := srv.Start(cfg.ServerAddress); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func buildHandlers(reg registry.UseCase, sessionRepo domainrepo.SessionRepository) *handlers {
	return &handlers{
		health:         publicHandler.NewHealthHandler(),
		fisherman:      adminHandler.NewFishermanHandler(reg),
		buyerAuth:      publicHandler.NewBuyerAuthHandler(reg, sessionRepo),
		buyer:          buyerHandler.NewBuyerHandler(reg),
		adminBuyer:     adminHandler.NewBuyerHandler(reg),
		publicItem:     publicHandler.NewItemHandler(reg),
		adminItem:      adminHandler.NewItemHandler(reg),
		bid:            buyerHandler.NewBidHandler(reg),
		invoice:        adminHandler.NewInvoiceHandler(reg),
		adminAuth:      publicHandler.NewAdminAuthHandler(reg, sessionRepo),
		publicVenue:    publicHandler.NewVenueHandler(reg),
		adminVenue:     adminHandler.NewVenueHandler(reg),
		publicAuction:  publicHandler.NewAuctionHandler(reg),
		adminAuction:   adminHandler.NewAuctionHandler(reg),
		admin:          adminHandler.NewAdminHandler(reg),
		authReset:      publicHandler.NewAuthResetHandler(reg),
		adminAuthReset: adminHandler.NewAuthResetHandler(reg),
		push:           buyerHandler.NewPushHandler(reg),
	}
}
