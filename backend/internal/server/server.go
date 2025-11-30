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

	"github.com/seka/fish-auction/backend/internal/server/handler"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
)

type Server struct {
	router           *http.ServeMux
	httpServer       *http.Server
	healthHandler    *handler.HealthHandler
	fishermanHandler *handler.FishermanHandler
	buyerHandler     *handler.BuyerHandler
	itemHandler      *handler.ItemHandler
	bidHandler       *handler.BidHandler
	invoiceHandler   *handler.InvoiceHandler
	authHandler      *handler.AuthHandler
	buyerAuth        *middleware.BuyerAuthMiddleware
}

func NewServer(
	healthHandler *handler.HealthHandler,
	fishermanHandler *handler.FishermanHandler,
	buyerHandler *handler.BuyerHandler,
	itemHandler *handler.ItemHandler,
	bidHandler *handler.BidHandler,
	invoiceHandler *handler.InvoiceHandler,
	authHandler *handler.AuthHandler,
) *Server {
	s := &Server{
		router:           http.NewServeMux(),
		healthHandler:    healthHandler,
		fishermanHandler: fishermanHandler,
		buyerHandler:     buyerHandler,
		itemHandler:      itemHandler,
		bidHandler:       bidHandler,
		invoiceHandler:   invoiceHandler,
		authHandler:      authHandler,
		buyerAuth:        middleware.NewBuyerAuthMiddleware(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.healthHandler.RegisterRoutes(s.router)
	s.fishermanHandler.RegisterRoutes(s.router)
	s.buyerHandler.RegisterRoutes(s.router)
	s.itemHandler.RegisterRoutes(s.router)

	// Protect bid routes with BuyerAuthMiddleware
	bidMux := http.NewServeMux()
	s.bidHandler.RegisterRoutes(bidMux)
	s.router.Handle("/api/bids", s.buyerAuth.Handle(bidMux))

	s.invoiceHandler.RegisterRoutes(s.router)
	s.authHandler.RegisterRoutes(s.router)
}

func (s *Server) Start(addr string) error {
	if addr == "" {
		addr = ":8080"
	}
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
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
