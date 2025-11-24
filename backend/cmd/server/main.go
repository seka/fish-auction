package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/seka/fish-auction/backend/internal/handler"
)

func main() {
	// Database Connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var db *sql.DB
	var err error

	// Retry connection
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("Failed to connect to DB: %v. Retrying in 2s...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	defer db.Close()

	// Run Migrations (Simple)
	// In a real app, use a migration tool. Here we just ensure tables exist.
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS fishermen (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS buyers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS auction_items (
		id SERIAL PRIMARY KEY,
		fisherman_id INTEGER REFERENCES fishermen(id),
		fish_type VARCHAR(255) NOT NULL,
		quantity INTEGER NOT NULL,
		unit VARCHAR(50) NOT NULL,
		status VARCHAR(50) DEFAULT 'Pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		item_id INTEGER REFERENCES auction_items(id),
		buyer_id INTEGER REFERENCES buyers(id),
		price INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(migrationSQL)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	h := handler.NewHandler(db)

	// Routes
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Backend is healthy!")
	})

	http.HandleFunc("/api/fishermen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateFisherman(w, r)
		} else if r.Method == http.MethodGet {
			h.GetFishermen(w, r)
		}
	})

	http.HandleFunc("/api/buyers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateBuyer(w, r)
		} else if r.Method == http.MethodGet {
			h.GetBuyers(w, r)
		}
	})

	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateItem(w, r)
		} else if r.Method == http.MethodGet {
			h.GetItems(w, r)
		}
	})

	http.HandleFunc("/api/bid", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.BidItem(w, r)
		}
	})

	http.HandleFunc("/api/invoices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetInvoices(w, r)
		}
	})

	http.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Login(w, r)
		}
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
