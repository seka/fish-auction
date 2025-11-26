package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	initializer "github.com/seka/fish-auction/backend/init"
	"github.com/seka/fish-auction/backend/internal/infrastructure/postgres"
	"github.com/seka/fish-auction/backend/internal/server"
	"github.com/seka/fish-auction/backend/internal/server/handler"
	"github.com/seka/fish-auction/backend/internal/usecase/auth"
	"github.com/seka/fish-auction/backend/internal/usecase/bid"
	"github.com/seka/fish-auction/backend/internal/usecase/buyer"
	"github.com/seka/fish-auction/backend/internal/usecase/fisherman"
	"github.com/seka/fish-auction/backend/internal/usecase/invoice"
	"github.com/seka/fish-auction/backend/internal/usecase/item"
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
	createFishermanUC := fisherman.NewCreateFishermanUseCase(fishermanRepo)
	listFishermenUC := fisherman.NewListFishermenUseCase(fishermanRepo)

	createBuyerUC := buyer.NewCreateBuyerUseCase(buyerRepo)
	listBuyersUC := buyer.NewListBuyersUseCase(buyerRepo)

	createItemUC := item.NewCreateItemUseCase(itemRepo)
	listItemsUC := item.NewListItemsUseCase(itemRepo)

	createBidUC := bid.NewCreateBidUseCase(itemRepo, bidRepo, txMgr)

	listInvoicesUC := invoice.NewListInvoicesUseCase(bidRepo)

	loginUC := auth.NewLoginUseCase()

	// Initialize Handlers
	healthHandler := handler.NewHealthHandler()
	fishermanHandler := handler.NewFishermanHandler(createFishermanUC, listFishermenUC)
	buyerHandler := handler.NewBuyerHandler(createBuyerUC, listBuyersUC)
	itemHandler := handler.NewItemHandler(createItemUC, listItemsUC)
	bidHandler := handler.NewBidHandler(createBidUC)
	invoiceHandler := handler.NewInvoiceHandler(listInvoicesUC)
	authHandler := handler.NewAuthHandler(loginUC)

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
