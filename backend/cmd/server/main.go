package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	initializer "github.com/seka/fish-auction/backend/init"
	"github.com/seka/fish-auction/backend/internal/domain/repository"
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

	// Initialize Dependencies
	repos := buildRepositories(db)
	txMgr := buildTransactionManager(db)
	useCases := buildUseCases(repos, txMgr)

	// Initialize Handlers
	handlers := buildHandlers(useCases)

	// Start Server
	srv := server.New(
		handlers.health,
		handlers.fisherman,
		handlers.buyer,
		handlers.item,
		handlers.bid,
		handlers.invoice,
		handlers.auth,
	)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

type repositories struct {
	fisherman repository.FishermanRepository
	buyer     repository.BuyerRepository
	item      repository.ItemRepository
	bid       repository.BidRepository
}

func buildRepositories(db *sql.DB) *repositories {
	return &repositories{
		fisherman: postgres.NewFishermanRepository(db),
		buyer:     postgres.NewBuyerRepository(db),
		item:      postgres.NewItemRepository(db),
		bid:       postgres.NewBidRepository(db),
	}
}

func buildTransactionManager(db *sql.DB) repository.TransactionManager {
	return postgres.NewTransactionManager(db)
}

type useCases struct {
	createFisherman fisherman.CreateFishermanUseCase
	listFishermen   fisherman.ListFishermenUseCase
	createBuyer     buyer.CreateBuyerUseCase
	listBuyers      buyer.ListBuyersUseCase
	createItem      item.CreateItemUseCase
	listItems       item.ListItemsUseCase
	createBid       bid.CreateBidUseCase
	listInvoices    invoice.ListInvoicesUseCase
	login           auth.LoginUseCase
}

func buildUseCases(repos *repositories, txMgr repository.TransactionManager) *useCases {
	return &useCases{
		createFisherman: fisherman.NewCreateFishermanUseCase(repos.fisherman),
		listFishermen:   fisherman.NewListFishermenUseCase(repos.fisherman),
		createBuyer:     buyer.NewCreateBuyerUseCase(repos.buyer),
		listBuyers:      buyer.NewListBuyersUseCase(repos.buyer),
		createItem:      item.NewCreateItemUseCase(repos.item),
		listItems:       item.NewListItemsUseCase(repos.item),
		createBid:       bid.NewCreateBidUseCase(repos.item, repos.bid, txMgr),
		listInvoices:    invoice.NewListInvoicesUseCase(repos.bid),
		login:           auth.NewLoginUseCase(),
	}
}

type handlers struct {
	health    *handler.HealthHandler
	fisherman *handler.FishermanHandler
	buyer     *handler.BuyerHandler
	item      *handler.ItemHandler
	bid       *handler.BidHandler
	invoice   *handler.InvoiceHandler
	auth      *handler.AuthHandler
}

func buildHandlers(ucs *useCases) *handlers {
	return &handlers{
		health:    handler.NewHealthHandler(),
		fisherman: handler.NewFishermanHandler(ucs.createFisherman, ucs.listFishermen),
		buyer:     handler.NewBuyerHandler(ucs.createBuyer, ucs.listBuyers),
		item:      handler.NewItemHandler(ucs.createItem, ucs.listItems),
		bid:       handler.NewBidHandler(ucs.createBid, ucs.listInvoices),
		invoice:   handler.NewInvoiceHandler(ucs.listInvoices),
		auth:      handler.NewAuthHandler(ucs.login),
	}
}
