package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/fish-auction/backend/internal/server/dto"
	"github.com/seka/fish-auction/backend/internal/server/util"
	"github.com/seka/fish-auction/backend/internal/usecase"
)

type InvoiceHandler struct {
	useCase usecase.InvoiceUseCase
}

func NewInvoiceHandler(uc usecase.InvoiceUseCase) *InvoiceHandler {
	return &InvoiceHandler{useCase: uc}
}

func (h *InvoiceHandler) List(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.useCase.List(r.Context())
	if err != nil {
		util.HandleError(w, err)
		return
	}

	resp := make([]dto.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		resp[i] = dto.InvoiceResponse{
			BuyerID:     inv.BuyerID,
			BuyerName:   inv.BuyerName,
			TotalAmount: inv.TotalAmount,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
