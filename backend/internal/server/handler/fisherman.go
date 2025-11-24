package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/model"
)

type FishermanHandler struct {
	db *sql.DB
}

func NewFishermanHandler(db *sql.DB) *FishermanHandler {
	return &FishermanHandler{db: db}
}

func (h *FishermanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var f model.Fisherman
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.QueryRow("INSERT INTO fishermen (name) VALUES ($1) RETURNING id", f.Name).Scan(&f.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

func (h *FishermanHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name FROM fishermen")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var fishermen []model.Fisherman
	for rows.Next() {
		var f model.Fisherman
		if err := rows.Scan(&f.ID, &f.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fishermen = append(fishermen, f)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fishermen)
}

func (h *FishermanHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/fishermen", func(w http.ResponseWriter, r *http.Request) {
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
