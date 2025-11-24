package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/model"
)

type InvoiceHandler struct {
	db *sql.DB
}

func NewInvoiceHandler(db *sql.DB) *InvoiceHandler {
	return &InvoiceHandler{db: db}
}

func (h *InvoiceHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
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

func (h *InvoiceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/invoices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.List(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
