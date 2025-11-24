package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/config"
	initializer "github.com/seka/fish-auction/backend/init"
	"github.com/seka/fish-auction/backend/internal/server"
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

	// Start Server
	srv := server.New(db)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
