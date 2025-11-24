package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/model"
)

type BidHandler struct {
	db *sql.DB
}

func NewBidHandler(db *sql.DB) *BidHandler {
	return &BidHandler{db: db}
}

func (h *BidHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update Item Status
	_, err = tx.Exec("UPDATE auction_items SET status = 'Sold' WHERE id = $1", t.ItemID)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert Transaction
	err = tx.QueryRow("INSERT INTO transactions (item_id, buyer_id, price) VALUES ($1, $2, $3) RETURNING id, created_at", t.ItemID, t.BuyerID, t.Price).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *BidHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/bid", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.Create(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
