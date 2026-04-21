package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/seka/fish-auction/backend/internal/domain/repository"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin"
	"github.com/seka/fish-auction/backend/internal/server/handler/buyer"
	"github.com/seka/fish-auction/backend/internal/server/handler/public"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
)

const (
	defaultShutdownTimeout = 30 * time.Second
)

// Server serves the request.
type Server struct {
	router                *http.ServeMux
	httpServer            *http.Server
	healthHandler         *public.HealthHandler
	fishermanHandler      *admin.FishermanHandler
	buyerAuthHandler      *public.BuyerAuthHandler
	buyerHandler          *buyer.BuyerHandler
	adminBuyerHandler     *admin.BuyerHandler
	publicItemHandler     *public.ItemHandler
	adminItemHandler      *admin.ItemHandler
	bidHandler            *buyer.BidHandler
	invoiceHandler        *admin.InvoiceHandler
	adminAuthHandler      *public.AdminAuthHandler
	publicVenueHandler    *public.VenueHandler
	adminVenueHandler     *admin.VenueHandler
	publicAuctionHandler  *public.AuctionHandler
	adminAuctionHandler   *admin.AuctionHandler
	adminHandler          *admin.AdminHandler
	adminAuthResetHandler *admin.AuthResetHandler
	authResetHandler      *public.AuthResetHandler
	pushHandler           *buyer.PushHandler
	adminAuth             *middleware.AdminAuthMiddleware
	buyerAuth             *middleware.BuyerAuthMiddleware
	cors                  *middleware.CORSMiddleware
	securityHeaders       *middleware.SecurityHeadersMiddleware
	csrf                  *middleware.CSRFMiddleware
	cacheControl          *middleware.CacheControlMiddleware
	gzip                  *middleware.GzipMiddleware
	maxBody               *middleware.MaxBodyMiddleware
	recovery              *middleware.RecoveryMiddleware
	readTimeout           time.Duration
	writeTimeout          time.Duration
	idleTimeout           time.Duration
}

// NewServer creates a new Server instance.
func NewServer(
	healthHandler *public.HealthHandler,
	fishermanHandler *admin.FishermanHandler,
	buyerAuthHandler *public.BuyerAuthHandler,
	buyerHandler *buyer.BuyerHandler,
	adminBuyerHandler *admin.BuyerHandler,
	publicItemHandler *public.ItemHandler,
	adminItemHandler *admin.ItemHandler,
	bidHandler *buyer.BidHandler,
	invoiceHandler *admin.InvoiceHandler,
	adminAuthHandler *public.AdminAuthHandler,
	publicVenueHandler *public.VenueHandler,
	adminVenueHandler *admin.VenueHandler,
	publicAuctionHandler *public.AuctionHandler,
	adminAuctionHandler *admin.AuctionHandler,
	adminHandler *admin.AdminHandler,
	authResetHandler *public.AuthResetHandler,
	adminAuthResetHandler *admin.AuthResetHandler,
	pushHandler *buyer.PushHandler,
	sessionRepo repository.SessionRepository,
	allowedOrigins []string,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	idleTimeout time.Duration,
) *Server {
	s := &Server{
		router:                http.NewServeMux(),
		healthHandler:         healthHandler,
		fishermanHandler:      fishermanHandler,
		buyerAuthHandler:      buyerAuthHandler,
		buyerHandler:          buyerHandler,
		adminBuyerHandler:     adminBuyerHandler,
		publicItemHandler:     publicItemHandler,
		adminItemHandler:      adminItemHandler,
		bidHandler:            bidHandler,
		invoiceHandler:        invoiceHandler,
		adminAuthHandler:      adminAuthHandler,
		publicVenueHandler:    publicVenueHandler,
		adminVenueHandler:     adminVenueHandler,
		publicAuctionHandler:  publicAuctionHandler,
		adminAuctionHandler:   adminAuctionHandler,
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
		readTimeout:           readTimeout,
		writeTimeout:          writeTimeout,
		idleTimeout:           idleTimeout,
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
	s.adminAuthHandler.RegisterRoutes(s.router)
	s.authResetHandler.RegisterRoutes(s.router)
	s.adminAuthResetHandler.RegisterRoutes(s.router)
	s.publicItemHandler.RegisterRoutes(s.router)
	s.publicAuctionHandler.RegisterRoutes(s.router)
	s.publicVenueHandler.RegisterRoutes(s.router)
	s.buyerAuthHandler.RegisterRoutes(s.router)
}

func (s *Server) registerAdminRoutes() {
	adminMux := http.NewServeMux()

	s.adminItemHandler.RegisterRoutes(adminMux)
	s.fishermanHandler.RegisterRoutes(adminMux)
	s.adminBuyerHandler.RegisterRoutes(adminMux)
	s.adminAuctionHandler.RegisterRoutes(adminMux)
	s.adminVenueHandler.RegisterRoutes(adminMux)
	s.adminHandler.RegisterRoutes(adminMux)
	s.invoiceHandler.RegisterRoutes(adminMux)

	s.router.Handle("/api/admin/", s.adminAuth.Handle(http.StripPrefix("/api/admin", adminMux)))
}

func (s *Server) registerBuyerRoutes() {
	buyerMux := http.NewServeMux()

	s.bidHandler.RegisterRoutes(buyerMux)
	s.buyerHandler.RegisterRoutes(buyerMux)
	s.pushHandler.RegisterRoutes(buyerMux)

	s.router.Handle("/api/buyer/", s.buyerAuth.Handle(http.StripPrefix("/api/buyer", buyerMux)))
}

// Start starts the HTTP server and blocks until the context is canceled.
func (s *Server) Start(ctx context.Context, addr string) error {
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
		ReadTimeout:       s.readTimeout,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      s.writeTimeout,
		IdleTimeout:       s.idleTimeout,
	}

	go func() {
		log.Printf("Server starting on %s\n", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
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
