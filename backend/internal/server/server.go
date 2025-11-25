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
)

type Server struct {
	router           *http.ServeMux
	healthHandler    *handler.HealthHandler
	fishermanHandler *handler.FishermanHandler
	buyerHandler     *handler.BuyerHandler
	itemHandler      *handler.ItemHandler
	bidHandler       *handler.BidHandler
	invoiceHandler   *handler.InvoiceHandler
	authHandler      *handler.AuthHandler
}

func New(
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
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.healthHandler.RegisterRoutes(s.router)
	s.fishermanHandler.RegisterRoutes(s.router)
	s.buyerHandler.RegisterRoutes(s.router)
	s.itemHandler.RegisterRoutes(s.router)
	s.bidHandler.RegisterRoutes(s.router)
	s.invoiceHandler.RegisterRoutes(s.router)
	s.authHandler.RegisterRoutes(s.router)
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exiting")
	return nil
}
