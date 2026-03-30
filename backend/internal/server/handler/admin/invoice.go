package admin

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/registry"
	"github.com/seka/fish-auction/backend/internal/server/handler/admin/response"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase/invoice"
)

// InvoiceHandler handles admin HTTP requests related to invoices.
type InvoiceHandler struct {
	listUseCase invoice.ListInvoicesUseCase
}

// NewInvoiceHandler creates a new InvoiceHandler instance.
func NewInvoiceHandler(r registry.UseCase) *InvoiceHandler {
	return &InvoiceHandler{
		listUseCase: r.NewListInvoicesUseCase(),
	}
}

// List handles the request to list invoices.
func (h *InvoiceHandler) List(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]response.Invoice, len(invoices))
	for i, inv := range invoices {
		resp[i] = response.Invoice{
			BuyerID:     inv.BuyerID,
			BuyerName:   inv.BuyerName,
			TotalAmount: inv.TotalAmount,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// RegisterRoutes registers the admin invoice handler routes to the given mux.
func (h *InvoiceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /invoices", h.List)
}
