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

	"github.com/gorilla/mux"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	"github.com/seka/fish-auction/backend/internal/server/middleware"
)

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
}

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
		adminAuth:             middleware.NewAdminAuthMiddleware(),
		buyerAuth:             middleware.NewBuyerAuthMiddleware(),
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
	// Health
	s.healthHandler.RegisterRoutes(s.router)

	// Auth (Login/Logout)
	s.router.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.authHandler.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	s.router.HandleFunc("/api/admin/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.authHandler.Logout(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	s.authResetHandler.RegisterRoutes(s.router)
	s.adminAuthResetHandler.RegisterRoutes(s.router)

	// Users Registration (List Only Publicly)
	s.router.HandleFunc("/api/fishermen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.fishermanHandler.List(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	s.router.HandleFunc("/api/buyers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.buyerHandler.List(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	s.router.HandleFunc("/api/buyers/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.buyerHandler.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Public Resources (Read Only)
	s.router.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.itemHandler.List(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	s.router.HandleFunc("/api/auctions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.auctionHandler.List(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	s.router.HandleFunc("/api/auctions/", func(w http.ResponseWriter, r *http.Request) {
		// Public Get Detail
		if r.Method == http.MethodGet {
			s.auctionHandler.Get(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	s.router.HandleFunc("/api/venues", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.venueHandler.List(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	s.router.HandleFunc("/api/venues/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.venueHandler.Get(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Legacy Invoice
	s.invoiceHandler.RegisterRoutes(s.router)
}

func (s *Server) registerAdminRoutes() {
	// Mount Admin Router
	adminRouter := mux.NewRouter()
	adminSub := adminRouter.PathPrefix("/api/admin").Subrouter()
	adminSub.Use(s.adminAuth.Handle)

	// Since we are moving to gorilla/mux, let's register the item routes properly
	// Note: ItemHandler.RegisterRoutes adds "/api/items" prefix,
	// but adminSub is already at "/api/admin".
	// The paths in RegisterRoutes might need adjustment or we register manually.

	// Items (Create, Update, Delete, Sort)
	adminSub.HandleFunc("/items", s.itemHandler.Create).Methods(http.MethodPost)
	adminSub.HandleFunc("/items/{id:[0-9]+}", s.itemHandler.Update).Methods(http.MethodPut)
	adminSub.HandleFunc("/items/{id:[0-9]+}", s.itemHandler.Delete).Methods(http.MethodDelete)
	adminSub.HandleFunc("/items/{id:[0-9]+}/sort-order", s.itemHandler.UpdateSortOrder).Methods(http.MethodPut)

	// Fishermen
	adminSub.HandleFunc("/fishermen", s.fishermanHandler.Create).Methods(http.MethodPost)
	adminSub.HandleFunc("/fishermen/{id:[0-9]+}", s.fishermanHandler.Delete).Methods(http.MethodDelete)

	// Buyers
	adminSub.HandleFunc("/buyers", s.buyerHandler.Create).Methods(http.MethodPost)
	adminSub.HandleFunc("/buyers/{id:[0-9]+}", s.buyerHandler.Delete).Methods(http.MethodDelete)

	// Auctions
	adminSub.HandleFunc("/auctions", s.auctionHandler.Create).Methods(http.MethodPost)
	adminSub.HandleFunc("/auctions/{id:[0-9]+}", s.auctionHandler.Update).Methods(http.MethodPut)
	adminSub.HandleFunc("/auctions/{id:[0-9]+}", s.auctionHandler.Delete).Methods(http.MethodDelete)
	adminSub.HandleFunc("/auctions/{id:[0-9]+}/status", s.auctionHandler.UpdateStatus).Methods(http.MethodPatch)
	adminSub.HandleFunc("/auctions/{id:[0-9]+}/reorder", s.itemHandler.Reorder).Methods(http.MethodPut)

	// Venues
	adminSub.HandleFunc("/venues", s.venueHandler.Create).Methods(http.MethodPost)
	adminSub.HandleFunc("/venues/{id:[0-9]+}", s.venueHandler.Update).Methods(http.MethodPut)
	adminSub.HandleFunc("/venues/{id:[0-9]+}", s.venueHandler.Delete).Methods(http.MethodDelete)

	// Settings
	adminSub.HandleFunc("/password", s.adminHandler.UpdatePassword).Methods(http.MethodPut)

	// Use the router for all /api/admin requests
	s.router.Handle("/api/admin/", adminRouter)
}

func (s *Server) registerBuyerRoutes() {
	buyerMux := http.NewServeMux()

	// Bids
	buyerMux.HandleFunc("/bids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.bidHandler.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// My Page related
	buyerMux.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.buyerHandler.GetCurrentBuyer(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	buyerMux.HandleFunc("/me/purchases", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.buyerHandler.GetMyPurchases(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	buyerMux.HandleFunc("/me/auctions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.buyerHandler.GetMyAuctions(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	buyerMux.HandleFunc("/password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			s.buyerHandler.UpdatePassword(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Push Notification
	buyerMux.HandleFunc("/push/subscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.pushHandler.Subscribe(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Mount Buyer Mux
	s.router.Handle("/api/buyer/", s.buyerAuth.Handle(http.StripPrefix("/api/buyer", buyerMux)))
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
