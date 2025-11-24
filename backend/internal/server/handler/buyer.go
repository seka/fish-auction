package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/model"
)

type BuyerHandler struct {
	db *sql.DB
}

func NewBuyerHandler(db *sql.DB) *BuyerHandler {
	return &BuyerHandler{db: db}
}

func (h *BuyerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var b model.Buyer
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.QueryRow("INSERT INTO buyers (name) VALUES ($1) RETURNING id", b.Name).Scan(&b.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

func (h *BuyerHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name FROM buyers")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var buyers []model.Buyer
	for rows.Next() {
		var b model.Buyer
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buyers = append(buyers, b)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buyers)
}

func (h *BuyerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/buyers", func(w http.ResponseWriter, r *http.Request) {
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
