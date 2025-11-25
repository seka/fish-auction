package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	initializer "github.com/seka/fish-auction/backend/init"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/seka/fish-auction/backend/internal/server"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	"github.com/seka/fish-auction/backend/internal/usecase"
)

func main() {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Database Connection
	connStr := cfg.DBConnectionURL()

	db, err := initializer.ConnectDB(connStr)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	defer db.Close()

	// Run Migrations
	if err := initializer.InitDB(db); err != nil {
		log.Fatal(err)
	}

	// Initialize Repositories
	fishermanRepo := postgres.NewFishermanRepository(db)
	buyerRepo := postgres.NewBuyerRepository(db)
	itemRepo := postgres.NewItemRepository(db)
	bidRepo := postgres.NewBidRepository(db)

	// Initialize Transaction Manager
	txMgr := postgres.NewTransactionManager(db)

	// Initialize Use Cases
	fishermanUseCase := usecase.NewFishermanInteractor(fishermanRepo)
	buyerUseCase := usecase.NewBuyerInteractor(buyerRepo)
	itemUseCase := usecase.NewItemInteractor(itemRepo)
	bidUseCase := usecase.NewBidInteractor(itemRepo, bidRepo, txMgr)
	invoiceUseCase := usecase.NewInvoiceInteractor(bidRepo)
	authUseCase := usecase.NewAuthInteractor()

	// Initialize Handlers
	healthHandler := handler.NewHealthHandler()
	fishermanHandler := handler.NewFishermanHandler(fishermanUseCase)
	buyerHandler := handler.NewBuyerHandler(buyerUseCase)
	itemHandler := handler.NewItemHandler(itemUseCase)
	bidHandler := handler.NewBidHandler(bidUseCase)
	invoiceHandler := handler.NewInvoiceHandler(invoiceUseCase)
	authHandler := handler.NewAuthHandler(authUseCase)

	// Start Server
	srv := server.New(
		healthHandler,
		fishermanHandler,
		buyerHandler,
		itemHandler,
		bidHandler,
		invoiceHandler,
		authHandler,
	)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
