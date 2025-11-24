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

	"github.com/seka/fish-auction/backend/internal/handler"
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
	h := handler.NewHandler(s.db)

	s.router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Backend is healthy!")
	})

	s.router.HandleFunc("/api/fishermen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateFisherman(w, r)
		} else if r.Method == http.MethodGet {
			h.GetFishermen(w, r)
		}
	})

	s.router.HandleFunc("/api/buyers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateBuyer(w, r)
		} else if r.Method == http.MethodGet {
			h.GetBuyers(w, r)
		}
	})

	s.router.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateItem(w, r)
		} else if r.Method == http.MethodGet {
			h.GetItems(w, r)
		}
	})

	s.router.HandleFunc("/api/bid", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.BidItem(w, r)
		}
	})

	s.router.HandleFunc("/api/invoices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetInvoices(w, r)
		}
	})

	s.router.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Login(w, r)
		}
	})
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
