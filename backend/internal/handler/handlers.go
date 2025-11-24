package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/model"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) CreateFisherman(w http.ResponseWriter, r *http.Request) {
	var f model.Fisherman
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.DB.QueryRow("INSERT INTO fishermen (name) VALUES ($1) RETURNING id", f.Name).Scan(&f.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

func (h *Handler) CreateBuyer(w http.ResponseWriter, r *http.Request) {
	var b model.Buyer
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.DB.QueryRow("INSERT INTO buyers (name) VALUES ($1) RETURNING id", b.Name).Scan(&b.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item model.AuctionItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.DB.QueryRow(
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

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	query := "SELECT id, fisherman_id, fish_type, quantity, unit, status, created_at FROM auction_items"
	var args []interface{}
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := h.DB.Query(query, args...)
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

func (h *Handler) BidItem(w http.ResponseWriter, r *http.Request) {
	var t model.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.DB.Begin()
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

func (h *Handler) GetInvoices(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT b.id, b.name, SUM(t.price) as total_price
		FROM transactions t
		JOIN buyers b ON t.buyer_id = b.id
		GROUP BY b.id, b.name
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var invoices []model.InvoiceItem
	for rows.Next() {
		var id int
		var name string
		var totalPrice int
		if err := rows.Scan(&id, &name, &totalPrice); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 8% Tax
		totalAmount := int(float64(totalPrice) * 1.08)

		invoices = append(invoices, model.InvoiceItem{
			BuyerID:     id,
			BuyerName:   name,
			TotalAmount: totalAmount,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if creds.Password == "admin-password" {
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_session",
			Value:    "authenticated",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   86400, // 1 day
		})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	} else {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
	}
}

func (h *Handler) GetFishermen(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name FROM fishermen")
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

func (h *Handler) GetBuyers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name FROM buyers")
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
