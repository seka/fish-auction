package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
)

// Server serves the request.
type Server struct {
	router                *http.ServeMux
	httpServer            *http.Server
	healthHandler         *handler.HealthHandler
	fishermanHandler      *handler.FishermanHandler
	buyerHandler          *handler.BuyerHandler
	itemHandler           *handler.ItemHandler
	bidHandler            *handler.BidHandler
	invoiceHandler        *handler.InvoiceHandler
	authHandler           *handler.AuthHandler
	venueHandler          *handler.VenueHandler
	auctionHandler        *handler.AuctionHandler
	adminHandler          *handler.AdminHandler
	adminAuthResetHandler *handler.AdminAuthResetHandler
	authResetHandler      *handler.AuthResetHandler
	pushHandler           *handler.PushHandler
	adminAuth             *middleware.AdminAuthMiddleware
	buyerAuth             *middleware.BuyerAuthMiddleware
	cors                  *middleware.CORSMiddleware
	securityHeaders       *middleware.SecurityHeadersMiddleware
	csrf                  *middleware.CSRFMiddleware
	cacheControl          *middleware.CacheControlMiddleware
	gzip                  *middleware.GzipMiddleware
	maxBody               *middleware.MaxBodyMiddleware
	recovery              *middleware.RecoveryMiddleware
	readTimeoutSec        time.Duration
	writeTimeoutSec       time.Duration
	idleTimeoutSec        time.Duration
}

// NewServer creates a new Server instance.
func NewServer(
	healthHandler *handler.HealthHandler,
	fishermanHandler *handler.FishermanHandler,
	buyerHandler *handler.BuyerHandler,
	itemHandler *handler.ItemHandler,
	bidHandler *handler.BidHandler,
	invoiceHandler *handler.InvoiceHandler,
	authHandler *handler.AuthHandler,
	venueHandler *handler.VenueHandler,
	auctionHandler *handler.AuctionHandler,
	adminHandler *handler.AdminHandler,
	authResetHandler *handler.AuthResetHandler,
	adminAuthResetHandler *handler.AdminAuthResetHandler,
	pushHandler *handler.PushHandler,
	sessionRepo repository.SessionRepository,
	allowedOrigins []string,
	readTimeoutSec time.Duration,
	writeTimeoutSec time.Duration,
	idleTimeoutSec time.Duration,
) *Server {
	s := &Server{
		router:                http.NewServeMux(),
		healthHandler:         healthHandler,
		fishermanHandler:      fishermanHandler,
		buyerHandler:          buyerHandler,
		itemHandler:           itemHandler,
		bidHandler:            bidHandler,
		invoiceHandler:        invoiceHandler,
		authHandler:           authHandler,
		venueHandler:          venueHandler,
		auctionHandler:        auctionHandler,
		adminHandler:          adminHandler,
		authResetHandler:      authResetHandler,
		adminAuthResetHandler: adminAuthResetHandler,
		pushHandler:           pushHandler,
		adminAuth:             middleware.NewAdminAuthMiddleware(sessionRepo),
		buyerAuth:             middleware.NewBuyerAuthMiddleware(sessionRepo),
		cors:                  middleware.NewCORSMiddleware(allowedOrigins),
		securityHeaders:       middleware.NewSecurityHeadersMiddleware(),
		csrf:                  middleware.NewCSRFMiddleware(allowedOrigins),
		cacheControl:          middleware.NewCacheControlMiddleware(),
		gzip:                  middleware.NewGzipMiddleware(),
		maxBody:               middleware.NewMaxBodyMiddleware(1024 * 1024), // 1MB
		recovery:              middleware.NewRecoveryMiddleware(),
		readTimeoutSec:        readTimeoutSec,
		writeTimeoutSec:       writeTimeoutSec,
		idleTimeoutSec:        idleTimeoutSec,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.registerPublicRoutes()
	s.registerAdminRoutes()
	s.registerBuyerRoutes()
}

func (s *Server) registerPublicRoutes() {
	s.healthHandler.RegisterRoutes(s.router)

	s.authHandler.RegisterRoutes(s.router)

	s.authResetHandler.RegisterRoutes(s.router)
	s.adminAuthResetHandler.RegisterRoutes(s.router)

	s.router.HandleFunc("GET /api/items", s.itemHandler.List)
	s.router.HandleFunc("GET /api/auctions", s.auctionHandler.List)
	s.router.HandleFunc("GET /api/auctions/{id}", s.auctionHandler.Get)
	s.router.HandleFunc("GET /api/auctions/{id}/items", s.auctionHandler.GetItems)
	s.router.HandleFunc("GET /api/venues", s.venueHandler.List)
	s.router.HandleFunc("GET /api/venues/{id}", s.venueHandler.Get)

	s.invoiceHandler.RegisterRoutes(s.router)
}

func (s *Server) registerAdminRoutes() {
	adminMux := http.NewServeMux()

	adminMux.HandleFunc("POST /items", s.itemHandler.Create)
	adminMux.HandleFunc("PUT /items/{id}", s.itemHandler.Update)
	adminMux.HandleFunc("DELETE /items/{id}", s.itemHandler.Delete)
	adminMux.HandleFunc("PUT /items/{id}/sort-order", s.itemHandler.UpdateSortOrder)

	adminMux.HandleFunc("GET /fishermen", s.fishermanHandler.List)
	adminMux.HandleFunc("POST /fishermen", s.fishermanHandler.Create)
	adminMux.HandleFunc("DELETE /fishermen/{id}", s.fishermanHandler.Delete)

	adminMux.HandleFunc("GET /buyers", s.buyerHandler.List)
	adminMux.HandleFunc("POST /buyers", s.buyerHandler.Create)
	adminMux.HandleFunc("DELETE /buyers/{id}", s.buyerHandler.Delete)

	adminMux.HandleFunc("POST /auctions", s.auctionHandler.Create)
	adminMux.HandleFunc("PUT /auctions/{id}", s.auctionHandler.Update)
	adminMux.HandleFunc("DELETE /auctions/{id}", s.auctionHandler.Delete)
	adminMux.HandleFunc("PATCH /auctions/{id}/status", s.auctionHandler.UpdateStatus)
	adminMux.HandleFunc("PUT /auctions/{id}/reorder", s.itemHandler.Reorder)

	adminMux.HandleFunc("POST /venues", s.venueHandler.Create)
	adminMux.HandleFunc("PUT /venues/{id}", s.venueHandler.Update)
	adminMux.HandleFunc("DELETE /venues/{id}", s.venueHandler.Delete)

	adminMux.HandleFunc("PUT /password", s.adminHandler.UpdatePassword)

	s.router.Handle("/api/admin/", s.adminAuth.Handle(http.StripPrefix("/api/admin", adminMux)))
}

func (s *Server) registerBuyerRoutes() {
	buyerMux := http.NewServeMux()

	buyerMux.HandleFunc("POST /bids", s.bidHandler.Create)

	buyerMux.HandleFunc("GET /me", s.buyerHandler.GetCurrentBuyer)
	buyerMux.HandleFunc("GET /me/purchases", s.buyerHandler.GetMyPurchases)
	buyerMux.HandleFunc("GET /me/auctions", s.buyerHandler.GetMyAuctions)
	buyerMux.HandleFunc("PUT /password", s.buyerHandler.UpdatePassword)

	buyerMux.HandleFunc("POST /push/subscribe", s.pushHandler.Subscribe)

	s.router.Handle("/api/buyer/", s.buyerAuth.Handle(http.StripPrefix("/api/buyer", buyerMux)))
}

// Start provides Start related functionality.
func (s *Server) Start(addr string) error {
	if addr == "" {
		addr = ":8080"
	}

	handlerWithCORS := s.cors.Handle(s.router)
	handlerWithHeaders := s.securityHeaders.Handle(handlerWithCORS)
	handlerWithCSRF := s.csrf.Handle(handlerWithHeaders)
	handlerWithCache := s.cacheControl.Handle(handlerWithCSRF)
	handlerWithGzip := s.gzip.Handle(handlerWithCache)
	handlerWithLimit := s.maxBody.Handle(handlerWithGzip)
	handlerWithRecovery := s.recovery.Handle(handlerWithLimit)

	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           handlerWithRecovery,
		ReadTimeout:       s.readTimeoutSec,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      s.writeTimeoutSec,
		IdleTimeout:       s.idleTimeoutSec,
	}

	go func() {
		log.Printf("Server starting on %s\n", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exiting")
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}
