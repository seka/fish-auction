package server

import (
	"context"
	"database/sql"
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
	db     *sql.DB
	router *http.ServeMux
}

func New(db *sql.DB) *Server {
	s := &Server{
		db:     db,
		router: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	healthHandler := handler.NewHealthHandler()
	fishermanHandler := handler.NewFishermanHandler(s.db)
	buyerHandler := handler.NewBuyerHandler(s.db)
	itemHandler := handler.NewItemHandler(s.db)
	bidHandler := handler.NewBidHandler(s.db)
	invoiceHandler := handler.NewInvoiceHandler(s.db)
	authHandler := handler.NewAuthHandler()

	healthHandler.RegisterRoutes(s.router)
	fishermanHandler.RegisterRoutes(s.router)
	buyerHandler.RegisterRoutes(s.router)
	itemHandler.RegisterRoutes(s.router)
	bidHandler.RegisterRoutes(s.router)
	invoiceHandler.RegisterRoutes(s.router)
	authHandler.RegisterRoutes(s.router)
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
