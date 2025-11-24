package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/model"
)

type ItemHandler struct {
	db *sql.DB
}

func NewItemHandler(db *sql.DB) *ItemHandler {
	return &ItemHandler{db: db}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var item model.AuctionItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.QueryRow(
		"INSERT INTO auction_items (fisherman_id, fish_type, quantity, unit, status) VALUES ($1, $2, $3, $4, 'Pending') RETURNING id, created_at, status",
		item.FishermanID, item.FishType, item.Quantity, item.Unit,
	).Scan(&item.ID, &item.CreatedAt, &item.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	query := "SELECT id, fisherman_id, fish_type, quantity, unit, status, created_at FROM auction_items"
	var args []interface{}
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []model.AuctionItem
	for rows.Next() {
		var i model.AuctionItem
		if err := rows.Scan(&i.ID, &i.FishermanID, &i.FishType, &i.Quantity, &i.Unit, &i.Status, &i.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, i)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *ItemHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			h.List(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
